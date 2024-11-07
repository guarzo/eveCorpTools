// Initialize the pilot list and pilot assignment on DOMContentLoaded
document.addEventListener("DOMContentLoaded", () => {
    console.log("Initializing pilot list...");
    initPilotList();

    const addPilotForm = document.getElementById("addPilotForm");
    if (addPilotForm) {
        addPilotForm.addEventListener("submit", addPilot);
    } else {
        console.error("addPilotForm not found in DOM.");
    }
});

// Fetches pilot names from the server and populates the pilot list
async function fetchPilotNames() {
    try {
        const response = await fetch('/api/pilots');
        const data = await response.json();

        const container = document.getElementById("pilot-list-container");
        if (!container) {
            console.error("pilot-list-container not found.");
            return;
        }

        container.innerHTML = `
            <div id="pilotNamesContainer" class="draggable-container flex flex-wrap gap-4 justify-center p-4"
                 ondrop="drop(event)" ondragover="allowDrop(event)">
            </div>`;

        const pilotContainer = document.getElementById("pilotNamesContainer");
        data.forEach(name => addPilotToContainer(name, pilotContainer, true));
    } catch (error) {
        console.error("Error fetching pilot names:", error);
    }
}

// Initialize the pilot list by fetching names from the server
function initPilotList() {
    fetchPilotNames();
}

// Add a pilot to the given container with an optional remove button
function addPilotToContainer(name, container, showRemoveButton = false) {
    if (!isPilotInContainer(name, container)) {
        const box = document.createElement('div');
        box.className = 'draggable-box bg-teal-500 text-gray-100 px-4 py-2 rounded-lg shadow-lg cursor-move border border-teal-500 hover:bg-teal-400 transition-all duration-200 flex items-center justify-between';


        box.draggable = true;
        box.ondragstart = (event) => event.dataTransfer.setData("text/plain", name);

        const pilotName = document.createElement('span');
        pilotName.className = 'mr-2';
        pilotName.textContent = name;

        if (showRemoveButton) {
            const removeButton = document.createElement('button');
            removeButton.className = 'text-red-500 hover:text-red-700 ml-2';
            removeButton.innerHTML = '<i class="fas fa-trash-alt"></i>';
            removeButton.onclick = (event) => {
                event.stopPropagation();
                removePilot(name);
            };
            box.appendChild(pilotName);
            box.appendChild(removeButton);
        } else {
            box.appendChild(pilotName);
        }

        container.appendChild(box);
    }
}

// Check if a pilot is already in a container to avoid duplicates
function isPilotInContainer(name, container) {
    return Array.from(container.getElementsByClassName("draggable-box"))
        .some(box => box.textContent.trim() === name);
}

// Event handler for the add pilot form submission
async function addPilot(event) {
    event.preventDefault();
    const pilotNameInput = document.getElementById("pilotNameInput");
    const name = pilotNameInput.value.trim();

    if (!name) {
        Swal.fire({
            icon: 'warning',
            title: 'Empty Name',
            text: 'Please enter a pilot name.',
        });
        return;
    }

    try {
        const response = await fetch('/api/pilots', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ name })
        });

        if (!response.ok) throw new Error(`HTTP error! Status: ${response.status}`);

        Swal.fire({
            icon: 'success',
            title: 'Added!',
            text: `Pilot "${name}" has been added.`,
            timer: 1500,
            showConfirmButton: false
        });

        pilotNameInput.value = '';
        fetchPilotNames();
    } catch (error) {
        console.error("Error adding pilot:", error);
        Swal.fire({
            icon: 'error',
            title: 'Failed to Add Pilot',
            text: 'An error occurred while adding the pilot. Please try again.',
        });
    }
}

// Remove pilot with confirmation dialog and server interaction
async function removePilot(name) {
    const result = await Swal.fire({
        title: 'Are you sure?',
        text: `Do you want to remove pilot "${name}"?`,
        icon: 'warning',
        showCancelButton: true,
        confirmButtonColor: '#d33',
        cancelButtonColor: '#3085d6',
        confirmButtonText: 'Yes, remove!'
    });

    if (result.isConfirmed) {
        try {
            const response = await fetch(`/api/pilots/${name}`, { method: 'DELETE' });
            if (!response.ok) throw new Error(`HTTP error! Status: ${response.status}`);

            Swal.fire({
                icon: 'success',
                title: 'Removed!',
                text: `Pilot "${name}" has been removed.`,
                timer: 1500,
                showConfirmButton: false
            });

            fetchPilotNames();
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
