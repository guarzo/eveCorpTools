// Initialize the pilot assignment areas with drop functionality
function initPilotAssignment() {
    const container = document.getElementById("pilot-assignment-container");
    container.innerHTML = `
        <div class="grid-container grid grid-cols-2 gap-4">
            <div id="scout" class="grid-box bg-gray-800 p-4 text-center text-lg rounded-lg shadow-lg border border-yellow-400" 
                 ondrop="drop(event)" ondragover="allowDrop(event)">
                <div class="text-yellow-400 font-bold mb-2">Scout</div>
                <div class="pilot-assignment-area flex flex-wrap gap-2 justify-center"></div>
            </div>
            <div id="involved" class="grid-box bg-gray-800 p-4 text-center text-lg rounded-lg shadow-lg border border-yellow-400" 
                 ondrop="drop(event)" ondragover="allowDrop(event)">
                <div class="text-yellow-400 font-bold mb-2">Involved</div>
                <div class="pilot-assignment-area flex flex-wrap gap-2 justify-center"></div>
            </div>
        </div>`;

    const pilotNamesContainer = document.getElementById("pilotNamesContainer");
    if (pilotNamesContainer) {
        pilotNamesContainer.ondrop = drop;
        pilotNamesContainer.ondragover = allowDrop;
    }
}

// Allow drop event
function allowDrop(event) {
    event.preventDefault();
}

// Drop event to add pilots to target containers
function drop(event) {
    event.preventDefault();
    const pilotName = event.dataTransfer.getData("text/plain"); // Get pilot name from drag event
    const dropBox = event.target.closest('.grid-box') || document.getElementById("pilotNamesContainer");

    if (!isPilotInContainer(pilotName, dropBox)) {
        removePilotFromAllContainers(pilotName, true); // Move action, no confirmation needed
        const showRemoveButton = dropBox.id === "pilotNamesContainer"; // Show remove button only in main list
        addPilotToContainer(pilotName, dropBox, showRemoveButton);

        const totalBuyPrice = parseFloat(document.getElementById("lootInput")?.innerText.replace(/,/g, '') || '0');
        calculateValues(totalBuyPrice);
    }
}

// Function to move pilot back to the main list
function movePilotBack(event) {
    const pilotName = event.target.textContent.trim();
    removePilotFromAllContainers(pilotName, true); // Move action, no confirmation needed
    addPilotToContainer(pilotName, document.getElementById("pilotNamesContainer"), true);

    const totalBuyPrice = parseFloat(document.getElementById("lootInput")?.innerText.replace(/,/g, '') || '0');
    calculateValues(totalBuyPrice);
}

// Function to remove pilot from all containers (includes SweetAlert confirmation and server interaction)
async function removePilotFromAllContainers(pilot, isMove = false) {
    if (isMove) {
        // If it's a move action, just remove the pilot without confirmation or server call
        ["pilotNamesContainer", "scout", "involved"].forEach(containerId => {
            const container = document.getElementById(containerId);
            if (container) {
                removePilotFromGrid(pilot, container);
            }
        });
        return;
    }

    // For delete action, show confirmation
    const result = await Swal.fire({
        title: 'Are you sure?',
        text: `Do you want to remove pilot "${pilot}"?`,
        icon: 'warning',
        showCancelButton: true,
        confirmButtonColor: '#d33',
        cancelButtonColor: '#3085d6',
        confirmButtonText: 'Yes, remove!'
    });

    if (result.isConfirmed) {
        ["pilotNamesContainer", "scout", "involved"].forEach(containerId => {
            const container = document.getElementById(containerId);
            if (container) {
                removePilotFromGrid(pilot, container);
            }
        });

        // Server interaction for removal
        try {
            const response = await fetch(`/api/pilots/${pilot}`, { method: 'DELETE' });
            if (!response.ok) throw new Error(`HTTP error! Status: ${response.status}`);

            Swal.fire({
                icon: 'success',
                title: 'Removed!',
                text: `Pilot "${pilot}" has been removed.`,
                timer: 1500,
                showConfirmButton: false
            });

            fetchPilotNames(); // Refresh list
        } catch (error) {
            console.error("Error removing pilot:", error);
            Swal.fire({
                icon: 'error',
                title: 'Failed to Remove Pilot',
                text: 'An error occurred while removing the pilot. Please try again.',
            });
        }
    }
}

// Function to remove a specific pilot from a specified container
function removePilotFromGrid(pilot, container) {
    const boxes = container.getElementsByClassName("draggable-box");
    for (const box of boxes) {
        if (box.textContent.trim() === pilot) {
            box.remove();
            break;
        }
    }
}

// Check if a pilot is already in the specified container
function isPilotInContainer(pilot, container) {
    const boxes = container.getElementsByClassName("draggable-box");
    return Array.from(boxes).some(box => box.textContent.trim() === pilot);
}

// Add a pilot to the container, with an optional remove button
function addPilotToContainer(pilotName, container, showRemoveButton) {
    const pilotBox = document.createElement("div");
    pilotBox.className = "draggable-box bg-cyan-500 text-white px-4 py-2 rounded-lg shadow-lg cursor-move border border-yellow-400 hover:bg-cyan-600 transition-all duration-200 flex items-center justify-between";
    pilotBox.draggable = true;
    pilotBox.ondragstart = function(event) {
        event.dataTransfer.setData("text/plain", pilotName);
    };

    const nameSpan = document.createElement("span");
    nameSpan.className = "mr-2";
    nameSpan.textContent = pilotName;

    pilotBox.appendChild(nameSpan);

    if (showRemoveButton) {
        const removeButton = document.createElement("button");
        removeButton.className = "text-red-500 hover:text-red-700 focus:outline-none";
        removeButton.innerHTML = '<i class="fas fa-trash-alt"></i>';
        removeButton.onclick = function(event) {
            event.stopPropagation();
            removePilotFromAllContainers(pilotName); // This triggers the SweetAlert delete confirmation
        };
        pilotBox.appendChild(removeButton);
    }

    const assignmentArea = container.querySelector(".pilot-assignment-area");
    if (assignmentArea) {
        assignmentArea.appendChild(pilotBox);
    } else {
        container.appendChild(pilotBox); // For the main pilot list area
    }
}

// Initialize the pilot list and pilot assignment on DOMContentLoaded
document.addEventListener("DOMContentLoaded", () => {
    console.log("Initializing pilot assignment...");
    initPilotAssignment()
});
