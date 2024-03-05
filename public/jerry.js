let width = 30;
let height = 30;
const resolution = 2;

let currentText = "";
let currentElement = null;

// Hide the default cursor
document.body.style.cursor = "none";

const el = document.createElement("div");
el.classList.add("jerry");

const canvas = document.createElement("div");
canvas.classList.add("jerry-canvas");
canvas.width = width * resolution;
canvas.height = height * resolution;
canvas.style.width = width;
canvas.style.height = height;

el.appendChild(canvas);

// Text text ele
const textEl = document.createElement("div");
textEl.classList.add("jerry-text");
el.appendChild(textEl);

// const ctx = el.getContext("2d");
// ctx.fillStyle = "#FF0000";
// ctx.fillRect(0, 0, width * resolution, height * resolution);

document.body.appendChild(el);

document.addEventListener("mousemove", (event) => {
  const { clientX, clientY } = event;

  el.style.transform = `translate(${clientX - width / 2}px, ${
    clientY - height / 2
  }px)`;

  const atPointer = document.elementFromPoint(clientX, clientY);
  const element = atPointer && atPointer.closest("[data-jerry]");
  const text = element && element.dataset.jerry;

  if (currentText !== text) {
    currentText = text;
    showElementText(currentText);

    canvas.style.backgroundColor = text ? "#FFFFFF" : "";
    canvas.style.transform = text ? "scale(0.2)" : "";
    // canvas.style.borderRadius = text ? width + "px" : "";

    currentElement && currentElement.classList.remove("jerry-on-me");
    element && element.classList.add("jerry-on-me");
    currentElement = element;
  }
});

function showElementText(text) {
  text = text && text.trim();
  if (text) {
    textEl.innerHTML = text;
    textEl.classList.add("showing");
  } else {
    // textEl.innerHTML = "";
    textEl.classList.remove("showing");
  }
}
