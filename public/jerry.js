// window.addEventListener("load", () => {
new Jerry();
// });

function debounce(mainFunction, delay) {
  // Declare a variable called 'timer' to store the timer ID
  let timer;

  // Return an anonymous function that takes in any number of arguments
  return function (...args) {
    // Clear the previous timer to prevent the execution of 'mainFunction'
    clearTimeout(timer);

    // Set a new timer that will execute 'mainFunction' after the specified delay
    timer = setTimeout(() => {
      mainFunction(...args);
    }, delay);
  };
}

function Jerry() {
  let width = 30;
  let height = 30;
  const parent = document.createElement("div");
  const canvas = document.createElement("div");
  const textEl = document.createElement("div");
  let tweening = new SmoothMover(
    ({ x, y }) => (parent.style.transform = `translate(${x}px, ${y}px)`)
  );

  hideBrowserCursor();
  assembleElements();

  // Keep track of the current element being focused
  let currentElement = null;

  const debouncedCheckUnderMouse = checkUnderTheMouse; //debounce(checkUnderTheMouse, 50);

  document.addEventListener("mousemove", (event) => {
    moveToPointer(event);

    debouncedCheckUnderMouse(event);
  });

  function checkUnderTheMouse(event) {
    // Find out what the element directly under the mouse is
    let underMouse = document.elementFromPoint(event.clientX, event.clientY);
    underMouse = underMouse && underMouse.closest("[data-jerry]");

    // Only update when pointing to something new
    if (currentElement === underMouse) return;

    setText(underMouse);
    setPointerStyle(underMouse);
    currentElement = underMouse;
  }

  function assembleElements() {
    parent.classList.add("jerry");
    canvas.classList.add("jerry-canvas");
    textEl.classList.add("jerry-text");

    parent.appendChild(canvas);
    parent.appendChild(textEl);
    // Add element to the DOM
    document.body.appendChild(parent);
  }

  function moveToPointer(event) {
    tweening.setTarget(event.clientX - width / 2, event.clientY - height / 2);
  }

  function setText(el) {
    let text = el && el.dataset.jerry;

    if (text) {
      textEl.innerHTML = text;
      textEl.classList.add("showing");
    } else {
      // Keep the text in there while it fades
      textEl.classList.remove("showing");
    }
  }

  function setPointerStyle(el) {
    if (el) {
      // Focus in
      canvas.style.backgroundColor = "#FFFFFF";
      canvas.style.transform = "scale(0.2)";
      return;
    }

    // Default style
    canvas.style.backgroundColor = "";
    canvas.style.transform = "";
  }
}

// We find related elements by checkign their dataset for a jerry field containing the text they which to display on the Jerry pointer
function taggedElementUnderMouse(event) {
  const el = document.elementFromPoint(event.clientX, event.clientY);
  return el && el.closest("[data-jerry]");
}

function hideBrowserCursor() {
  document.body.style.cursor = "none";
}

function SmoothMover(moveFn) {
  const position = { x: 0, y: 0 };
  let targetPosition = { x: 0, y: 0 };
  let running = true;
  requestAnimationFrame(update);

  function update() {
    position.x = position.x * 0.6 + targetPosition.x * 0.4;
    position.y = position.y * 0.6 + targetPosition.y * 0.4;
    const dx = Math.abs(position.x === targetPosition.x);
    const dy = Math.abs(position.y === targetPosition.y);
    const farFromDone = dx > 0.1 || dy > 0.1;
    if (farFromDone) {
      moveFn(position);
      requestAnimationFrame(update);
      return;
    }

    // Snap to the target and stop the animation once its close enough
    position.x = targetPosition.x;
    position.y = targetPosition.y;
    moveFn(position);

    running = false;
  }

  return {
    setTarget: (x, y) => {
      targetPosition.x = x;
      targetPosition.y = y;

      if (!running) {
        running = true;
        requestAnimationFrame(update);
      }
    },
  };
}
