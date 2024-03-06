Jerry();

function Jerry() {
  let width = 30;
  let height = 30;
  const parent = document.createElement("div");
  const canvas = document.createElement("div");
  const textEl = document.createElement("div");
  const position = { x: 0, y: 0 };
  const targetPosition = { x: 0, y: 0 };

  hideBrowserCursor();
  assembleElements();

  // Keep track of the current element being focused
  let currentElement = null;

  document.addEventListener("mousemove", (event) => {
    moveToPointer(event);

    // Find out what the element directly under the mouse is
    let underMouse = document.elementFromPoint(event.clientX, event.clientY);
    underMouse = underMouse && underMouse.closest("[data-jerry]");

    // Only update when pointing to something new
    if (currentElement === underMouse) return;

    setText(underMouse);
    setPointerStyle(underMouse);
    currentElement = underMouse;
  });

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
    targetPosition.x = event.clientX - width / 2;
    targetPosition.y = event.clientY - height / 2;

    tweening.setTarget(targetPosition);
  }

  let tweening = smoothMover();

  function smoothMover() {
    const position = { x: 0, y: 0 };
    let targetPosition = { x: 0, y: 0 };

    requestAnimationFrame(update);

    function update() {
      position.x = position.x * 0.6 + targetPosition.x * 0.4;
      position.y = position.y * 0.6 + targetPosition.y * 0.4;

      parent.style.transform = `translate(${position.x}px, ${position.y}px)`;

      if (
        !(position.x === targetPosition.x && position.y === targetPosition.y)
      ) {
        requestAnimationFrame(update); // Continue the animation if the duration has not been reached
      }
    }

    return {
      setTarget: (pos) => {
        targetPosition = pos;
      },
    };
  }

  function setText(el) {
    let text = el && el.dataset.jerry;
    text = text && text.trim();

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
