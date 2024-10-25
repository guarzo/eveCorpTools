function initPilotAssignment() {
    const container = document.getElementById("pilot-assignment-container");
    container.innerHTML = `
        <div class="grid-container">
            <div id="scout" class="grid-box" ondrop="drop(event)" ondragover="allowDrop(event)">Scout</div>
            <div id="involved" class="grid-box" ondrop="drop(event)" ondragover="allowDrop(event)">Involved</div>
        </div>`;

    document.getElementById("scout").addEventListener("dblclick", movePilotBack);
    document.getElementById("involved").addEventListener("dblclick", movePilotBack);
}

function drop(event) {
    event.preventDefault();
    const data = event.dataTransfer.getData("text");
    const dropBox = event.target.closest('.grid-box') || document.getElementById("characterNamesContainer");

    if (!isPilotInContainer(data, dropBox)) {
        // Remove pilot from all containers
        removePilotFromAllContainers(data);

        // Add pilot to the drop location
        addPilotToContainer(data, dropBox);

        const lootInputElement = document.getElementById("lootInput");
        const lootInputValue = lootInputElement ? lootInputElement.innerText.replace(/,/g, '') : '0'; // Handle undefined
        const totalBuyPrice = parseFloat(lootInputValue);

        calculateValues(totalBuyPrice);
    }
}

function movePilotBack(event) {
    const pilotName = event.target.innerHTML;
    removePilotFromAllContainers(pilotName);
    const characterNamesContainer = document.getElementById("characterNamesContainer");
    addPilotToContainer(pilotName, characterNamesContainer);

    const lootInputElement = document.getElementById("lootInput");
    const lootInputValue = lootInputElement ? lootInputElement.innerText.replace(/,/g, '') : '0'; // Handle undefined
    const totalBuyPrice = parseFloat(lootInputValue);

    calculateValues(totalBuyPrice);
}

function removePilotFromAllContainers(pilot) {
    const containers = [
        document.getElementById("characterNamesContainer"),
        document.getElementById("scout"),
        document.getElementById("involved")
    ];
    containers.forEach(container => removePilotFromGrid(pilot, container));
}

function removePilotFromGrid(pilot, gridBox) {
    const boxes = gridBox.getElementsByClassName("draggable-box");
    for (const box of boxes) {
        if (box.innerHTML === pilot) {
            box.remove();
            break;
        }
    }
}
