function initCharacterList() {
    fetchCharacterNames();
}

async function fetchCharacterNames() {
    try {
        const response = await fetch('/fetch-character-names');
        const data = await response.json();
        const container = document.getElementById("character-list-container");
        container.innerHTML = `
            <div id="characterNamesContainer" class="draggable-container" ondrop="drop(event)" ondragover="allowDrop(event)">
            </div>`;
        const characterContainer = document.getElementById("characterNamesContainer");
        characterContainer.innerHTML = '';
        data.forEach(name => {
            addPilotToContainer(name, characterContainer);
        });
    } catch (error) {
        console.error("Error fetching character names.", error);
    }
}

function addPilotToContainer(name, container) {
    if (!isPilotInContainer(name, container)) {
        const box = document.createElement('div');
        box.className = 'draggable-box';
        box.draggable = true;
        box.innerHTML = name;
        box.ondragstart = drag;
        container.appendChild(box);
    }
}

function isPilotInContainer(name, container) {
    const boxes = container.getElementsByClassName("draggable-box");
    for (const box of boxes) {
        if (box.innerHTML === name) {
            return true;
        }
    }
    return false;
}

function allowDrop(event) {
    event.preventDefault();
}

function drag(event) {
    event.dataTransfer.setData("text", event.target.innerHTML);
}
