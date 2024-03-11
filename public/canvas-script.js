const canvasContainer = document.body;

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

  function handleObserved(entries) {
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

function luminance([r, g, b]) {
  return (0.2126 * r) / 255 + (0.7152 * g) / 255 + (0.0722 * b) / 255;
}
