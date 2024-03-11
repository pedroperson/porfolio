const RESOLUTION = 1 / 1;

class Rect {
  constructor(x, y, width, height, endX, endY, startTime, duration, color) {
    this.x = x * RESOLUTION;
    this.y = y * RESOLUTION;
    this.width = width * RESOLUTION;
    this.height = height * RESOLUTION;
    this.startX = x * RESOLUTION;
    this.startY = y * RESOLUTION;
    this.endX = endX * RESOLUTION;
    this.endY = endY * RESOLUTION;
    this.startTime = startTime;
    this.duration = duration;
    this.done = false;
    this.color = color;
  }

  step(currentTime) {
    const r = (currentTime - this.startTime) / this.duration;
    // TODO: THis function is wrong in some cruciial way
    const p = cubicBezier(r, 1, 0);

    if (p >= 1) {
      if (!this.done) {
        this.done = true;
        this.x = this.endX;
        this.y = this.endY;
      }
      return;
    }

    this.x = this.startX * (1 - p) + this.endX * p;
    this.y = this.startY * (1 - p) + this.endY * p;
  }
}

// between 0 and 1
function cubicBezier(t, p1, p2) {
  return 3 * (1 - t) * (1 - t) * t * p1 + 3 * (1 - t) * t * t * p2 + t * t * t;
}

function now() {
  return performance.now();
}

function lighten([r, g, b], ratio) {
  const w = ratio > 0 ? 255 : 0;
  r += (w - r) * ratio;
  g += (w - g) * ratio;
  b += (w - b) * ratio;
  return `rgb(${r},${g},${b})`;
}

function rectFromLeft(y, h, color = "", delay = 0, duration = 500) {
  const w = window.innerWidth;
  return new Rect(
    -w,
    y + 40 * 2 * (Math.random() - 0.5),
    w,
    h,
    0,
    y,
    now() + delay,
    duration,
    color
  );
}
function rectFromRight(y, h, color = "", delay = 0, duration = 500) {
  const w = window.innerWidth;
  return new Rect(
    w,
    y + h * 0.5 * (Math.random() - 0.5),
    w,
    h,
    0,
    y,
    now() + delay,
    duration,
    color
  );
}

const RECTS = [];
function startFrameChange(color, N = 4) {
  //   Define the drawables
  const h = Math.ceil(window.innerHeight / N);
  for (let i = 0; i < N; i++) {
    RECTS.push(
      rectFromRight(i * h, h, lighten(color, 0.05 * (Math.random() - 0.5)))
    );
  }
}

// for (let i = 0; i < 100; i++) {
//   setTimeout(() => {
//     startFrameChange(
//       [255 * Math.random(), 255 * Math.random(), 255 * Math.random()],
//       7
//     );
//   }, i * 1000);
// }

//   Prepare the canvas
const canvas = document.createElement("canvas");
canvas.width = window.innerWidth * RESOLUTION;
canvas.height = window.innerHeight * RESOLUTION;

const canvasContainer = document.body;
canvasContainer.appendChild(canvas);

const ctx = canvas.getContext("2d");

const white = "#EEEEEE";
const black = "#111111";
let bodyTextColor = white;
const colorCache = [];
const textColorCache = [];
const elementCache = [];
let currentColor = null;
let currentTextColor = null;

// Setup the trigger zones
window.addEventListener("DOMContentLoaded", () => {
  let observer = new IntersectionObserver(handleObserved, {
    rootMargin: "-50% 0% -50% 0%",
  });

  Array.from(document.querySelectorAll(".js-bg-color")).forEach((e) =>
    observer.observe(e)
  );

  function handleObserved(entries, observer) {
    entries.forEach((entry) => {
      if (!entry.isIntersecting) return;
      let elem = entry.target;

      let cacheIndex = -1;
      for (let ci = 0; ci < elementCache.length; ci++) {
        if (elem === elementCache[ci]) {
          cacheIndex = ci;
          break;
        }
      }

      const color =
        cacheIndex !== -1 ? colorCache[cacheIndex] : readElementColor(elem);

      // startFrameChange(color, 5);
      let textColor =
        cacheIndex !== -1
          ? textColorCache[cacheIndex]
          : luminance(color) < 0.5
          ? white
          : black;

      if (cacheIndex === -1) {
        elementCache.push(elem);
        colorCache.push(color);
        textColorCache.push(textColor);
      }

      setColor(color, textColor);
    });
  }

  function setColor(color, textColor) {
    if (currentColor && currentColor.every((c, i) => c === color[i])) {
      return;
    }
    canvasContainer.style.background = `rgb(${color[0]},${color[1]},${color[2]})`;

    if (!currentTextColor || currentTextColor !== textColor) {
      canvasContainer.style.color = textColor;
    }

    currentColor = color;
    currentTextColor = textColor;
  }
});

function readElementColor(elem) {
  const data = elem.getAttribute("data-color");
  return data ? JSON.parse(data) : [255, 255, 255];
}

//   Start the animation
requestAnimationFrame(animate);

function animate() {
  const n = now();

  let lastestDone = -1;
  for (let i = 0; i < RECTS.length; i++) {
    const r = RECTS[i];
    if (r.done) {
      lastestDone = i;
      continue;
    }

    r.step(n);
    ctx.fillStyle = r.color;
    ctx.fillRect(r.x, r.y, r.width, r.height);
  }

  //   Remove bottom most layer if it is covered by a next layer
  if (lastestDone > 0) {
    const oldestTime = RECTS[0].startTime;
    const secondFinishedLayer = RECTS.findIndex(
      (r) => r.done && r.startTime !== oldestTime
    );
    if (secondFinishedLayer !== 1) {
      RECTS.splice(0, secondFinishedLayer);
    }
  }

  requestAnimationFrame(animate);
}

function luminance([r, g, b]) {
  return (0.2126 * r) / 255 + (0.7152 * g) / 255 + (0.0722 * b) / 255;
}
