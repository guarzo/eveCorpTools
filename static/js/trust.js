// trust.js

// Initialize state variables
let activeRequests = 0; // Counter for active requests
let isShowingUntrusted = false; // Default view

// Initialize table variables grouped under a single object
const tables = {};

// Configure Toastr options for notifications
toastr.options = {
    "closeButton": true,
    "debug": false,
    "newestOnTop": true,
    "progressBar": true,
    "positionClass": "toast-top-right",
    "preventDuplicates": false,
    "onclick": null,
    "showDuration": "300",
    "hideDuration": "1000",
    "timeOut": "1500", // Increased timeout for better visibility
    "extendedTimeOut": "1000",
    "showEasing": "swing",
    "hideEasing": "linear",
    "showMethod": "fadeIn",
    "hideMethod": "fadeOut"
};

console.log(TrustedCharacters)
console.log(TrustedCorporations)

/**
 * Function to show the loading indicator
 */
function showLoading() {
    activeRequests += 1;
    const loader = document.getElementById("loading-indicator");
    if (loader) {
        loader.classList.remove("hidden");
        console.log("showLoading called. activeRequests:", activeRequests);
    }
}

/**
 * Updates the "Write Contacts for All" button's disabled state based on character tiles.
 */
function updateWriteAllButtonState() {
    const writeAllBtn = document.getElementById("write-contacts-all-btn");
    const characterTiles = document.querySelectorAll(".character-tile");
    if (writeAllBtn) {
        writeAllBtn.disabled = characterTiles.length === 0;
        writeAllBtn.style.cursor = characterTiles.length === 0 ? "not-allowed" : "pointer";
        writeAllBtn.style.opacity = characterTiles.length === 0 ? "0.6" : "1";
    }
}


function hideLoading() {
    if (activeRequests > 0) {
        activeRequests -= 1;
    } else {
        console.warn("hideLoading called but activeRequests is already 0 or negative");
    }

    // Clamp activeRequests to 0 to prevent it from being negative
    if (activeRequests < 0) {
        activeRequests = 0;
        console.warn("activeRequests clamped to 0");
    }

    const loader = document.getElementById("loading-indicator");
    if (loader && activeRequests <= 0) {
        loader.classList.add("hidden");
        console.log("hideLoading called. activeRequests:", activeRequests);
    }
}

/**
 * Function to toggle the disabled state of a button
 * @param {HTMLElement} button - The button element to toggle
 * @param {boolean} disable - Whether to disable or enable the button
 */
function toggleButtonState(button, disable) {
    if (!button) return;

    button.disabled = disable;
    button.style.cursor = disable ? "not-allowed" : "pointer";
    button.style.opacity = disable ? "0.6" : "1";
}

/**
 * Function to resize a Tabulator table container
 * @param {string} tableId - The ID of the table container
 */
function resizeTabulatorTable(tableId) {
    const tableInstance = tables[tableId];
    if (tableInstance) {
        tableInstance.redraw(true);
    }
}

/**
 * General function to toggle visibility of a section
 * @param {string} sectionId - The ID of the section to toggle
 * @param {boolean} show - Whether to show or hide the section
 */
function toggleSectionDisplay(sectionId, show) {
    const section = document.getElementById(sectionId);
    if (!section) {
        console.error(`Section with ID "${sectionId}" not found.`);
        return;
    }
    section.style.visibility = show ? "visible" : "hidden";
    section.style.opacity = show ? "1" : "0";
    section.style.position = show ? "relative" : "absolute";
}

/**
 * Function to toggle multiple sections
 * @param {Array} sections - Array of section IDs to toggle
 * @param {boolean} show - Whether to show or hide the sections
 */
function toggleMultipleSections(sections, show) {
    sections.forEach(sectionId => toggleSectionDisplay(sectionId, show));
}

/**
 * Helper function to determine the correct server endpoint
 * @param {string} trustStatus - 'trusted' or 'untrusted'
 * @param {string} entityType - 'character' or 'corporation'
 * @param {string} action - 'add' or 'remove'
 * @returns {string|null} - The server endpoint URL or null if invalid inputs
 */
function getServerEndpoint(trustStatus, entityType, action) {
    const endpoints = {
        trusted: {
            character: {
                add: '/validate-and-add-trusted-character',
                remove: '/remove-trusted-character',
            },
            corporation: {
                add: '/validate-and-add-trusted-corporation',
                remove: '/remove-trusted-corporation',
            }
        },
        untrusted: {
            character: {
                add: '/validate-and-add-untrusted-character',
                remove: '/remove-untrusted-character',
            },
            corporation: {
                add: '/validate-and-add-untrusted-corporation',
                remove: '/remove-untrusted-corporation',
            }
        }
    };

    return endpoints[trustStatus]?.[entityType]?.[action] || null;
}

/**
 * Centralized fetch function that handles JSON and text responses.
 * Throws an error with the appropriate message based on response status.
 * @param {string} url - The endpoint URL.
 * @param {object} options - Fetch options.
 * @returns {Promise<Object|string>} - Parsed JSON object or plain text string.
 */
async function fetchWithHandling(url, options) {
    try {
        const response = await fetch(url, options);
        const contentType = response.headers.get("Content-Type");

        if (!response.ok) {
            let errorData;
            if (contentType && contentType.includes("application/json")) {
                errorData = await response.json();
                throw new Error(errorData.error || "An error occurred.");
            } else {
                errorData = await response.text();
                throw new Error(errorData || "An error occurred.");
            }
        }

        if (contentType && contentType.includes("application/json")) {
            return response.json();
        } else {
            return response.text();
        }
    } catch (error) {
        console.error(`Network or parsing error: ${error}`);
        throw error;
    }
}

/**
 * Updates the progress bar.
 * @param {number} total - Total number of operations.
 * @param {number} completed - Number of completed operations.
 */
function updateProgressBar(total, completed) {
    const progressBar = document.getElementById("progress-bar");
    if (progressBar) {
        const percentage = total === 0 ? 0 : (completed / total) * 100;
        progressBar.style.width = `${percentage}%`;
    }
}

/**
 * Resets the progress bar.
 */
function resetProgressBar() {
    const progressBar = document.getElementById("progress-bar");
    if (progressBar) {
        progressBar.style.width = `0%`;
    }
}


/**
 * Adds an event listener to the "Write Contacts for All" button.
 */
function setupWriteAllButton() {
    const writeAllBtn = document.getElementById("write-contacts-all-btn");
    if (writeAllBtn) {
        writeAllBtn.addEventListener("click", function () {
            console.log("Write Contacts for All button clicked.");
            writeContactsForAll();
        });
    } else {
        console.error(`Button with ID "write-contacts-all-btn" not found.`);
    }
}

async function writeContactsForAll() {
    const writeAllBtn = document.getElementById("write-contacts-all-btn");
    const progressBar = document.getElementById("write-all-progress-bar");
    if (!writeAllBtn || !progressBar) {
        console.error(`Button or progress bar not found.`);
        return;
    }

    // Disable the button to prevent multiple clicks
    toggleButtonState(writeAllBtn, true);
    showLoading();

    const totalCharacters = TabulatorIdentities.length;
    let completed = 0;
    let successCount = 0;
    let failureCount = 0;

    // Reset progress bar
    progressBar.style.width = '0%';

    // Function to update progress
    const updateProgress = () => {
        completed += 1;
        const percentage = Math.round((completed / totalCharacters) * 100);
        progressBar.style.width = `${percentage}%`;
    };

    // Array to hold promises
    const promises = TabulatorIdentities.map(character => {
        return writeContacts(character.CharacterID)
            .then(() => {
                successCount += 1;
            })
            .catch(() => {
                failureCount += 1;
            })
            .finally(() => {
                updateProgress();
            });
    });

    try {
        await Promise.all(promises);
        if (failureCount === 0) {
            toastr.success(`Successfully updated contacts for all ${totalCharacters} characters.`);
        } else if (successCount === 0) {
            toastr.error(`Failed to update contacts for all ${totalCharacters} characters.`);
        } else {
            toastr.warning(`Updated contacts for ${successCount} characters. Failed for ${failureCount} characters.`);
        }
    } catch (error) {
        console.error(`Error writing contacts for all characters: ${error}`);
        toastr.error("An unexpected error occurred while writing contacts for all characters.");
    } finally {
        // Reset the progress bar after completion
        progressBar.style.width = '0%'; // Reset to 0%

        // Optionally, add a slight delay for visual feedback
        setTimeout(() => {
            // If you want to hide the progress bar after resetting
            // You can add a class to hide it or manipulate its visibility
            // For example, removing the width might suffice if it's set to 0%
            // Or you can toggle a hidden class
            // Here, we'll ensure it's at 0% and visible for the next operation
        }, 300); // Adjust the timeout as needed based on transition duration

        hideLoading();
        toggleButtonState(writeAllBtn, false);
    }
}

/**
 * Function to write contacts
 * Calls /add-contacts and /delete-contacts endpoints sequentially
 * @param {number} characterID - ID of the character
 */
async function writeContacts(characterID) {
    // Validate characterID
    if (typeof characterID !== 'number' || isNaN(characterID) || characterID <= 0) {
        toastr.error("Invalid Character ID.");
        console.error("Invalid Character ID:", characterID);
        return;
    }

    showLoading();
    const toggleBtn = document.getElementById("toggle-contacts-btn");
    toggleButtonState(toggleBtn, true);

    try {
        console.log("Sending to /add-contacts:", { characterID });

        // Call /add-contacts endpoint
        let response = await fetch(`/add-contacts`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ characterID })
        });

        if (!response.ok) {
            let errorMessage = "Failed to add contacts.";
            try {
                const errData = await response.json();
                errorMessage = errData.error || errorMessage;
            } catch (e) {
                const errText = await response.text();
                errorMessage = errText || errorMessage;
            }
            throw new Error(errorMessage);
        }

        const addData = await response.json();
        console.log("Contacts added successfully:", addData);

        console.log("Sending to /delete-contacts:", { characterID });

        // Call /delete-contacts endpoint
        response = await fetch(`/delete-contacts`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ characterID })
        });

        if (!response.ok) {
            let errorMessage = "Failed to delete contacts.";
            try {
                const errData = await response.json();
                errorMessage = errData.error || errorMessage;
            } catch (e) {
                const errText = await response.text();
                errorMessage = errText || errorMessage;
            }
            throw new Error(errorMessage);
        }

        const deleteData = await response.json();
        console.log("Contacts deleted successfully:", deleteData);
        toastr.success("Contacts updated successfully.");
    } catch (error) {
        toastr.error("Error writing contacts: " + error.message);
        console.error("Error writing contacts:", error);
    } finally {
        hideLoading();
        toggleButtonState(toggleBtn, false);
        updateWriteAllButtonState();
    }
}

/**
 * Checks if an entity is in a given list by Identifier (ID or Name).
 * @param {Array} list - The list to check.
 * @param {string|number} identifier - The ID or Name of the entity.
 * @param {string} entityType - 'character' or 'corporation'.
 * @returns {boolean} - True if the entity is in the list.
 */
function isEntityInListByIdentifier(list, identifier, entityType) {
    if (!identifier) {
        console.warn(`Invalid identifier:`, identifier);
        return false;
    }

    const isId = typeof identifier === 'number' || /^\d+$/.test(identifier);
    const idField = entityType === 'character' ? 'CharacterID' : 'CorporationID';
    const nameField = entityType === 'character' ? 'CharacterName' : 'CorporationName';

    if (isId) {
        const numericId = typeof identifier === 'number' ? identifier : parseInt(identifier, 10);
        if (isNaN(numericId)) {
            console.warn(`Identifier "${identifier}" is not a valid number.`);
            return false;
        }
        return list.some(entity => entity[idField] === numericId);
    } else {
        if (typeof identifier !== 'string') {
            console.warn(`Identifier must be a string when not numeric.`);
            return false;
        }
        return list.some(entity => {
            const entityName = entity[nameField];
            if (typeof entityName !== 'string') {
                console.warn(`Invalid ${nameField} in entity:`, entity);
                return false;
            }
            return entityName.toLowerCase() === identifier.toLowerCase();
        });
    }
}

/** Helper functions for character and corporation checks */
const isCharacterInTrustedList = (identifier) => isEntityInListByIdentifier(TrustedCharacters, identifier, 'character');
const isCorporationInTrustedList = (identifier) => isEntityInListByIdentifier(TrustedCorporations, identifier, 'corporation');
const isCharacterInUntrustedList = (identifier) => isEntityInListByIdentifier(UntrustedCharacters, identifier, 'character');
const isCorporationInUntrustedList = (identifier) => isEntityInListByIdentifier(UntrustedCorporations, identifier, 'corporation');

/**
 * Checks if a character is trusted based on character and corporation trust lists.
 * @param {object} character - The character object.
 * @returns {boolean} - True if trusted, otherwise false.
 */
function isCharacterTrusted(character) {
    return isCharacterInTrustedListByID(character.CharacterID) || isCorporationInTrustedListByID(character.CorporationID);
}

const isCharacterInTrustedListByID = (id) => isEntityInListByIdentifier(TrustedCharacters, id, 'character');
const isCorporationInTrustedListByID = (id) => isEntityInListByIdentifier(TrustedCorporations, id, 'corporation');

/**
 * Updates the trust status of a character tile
 * @param {HTMLElement} tile - The character tile element
 * @param {boolean} isTrusted - Whether the character is trusted
 */
function updateTileTrustStatus(tile, isTrusted) {
    console.log(`Updating tile trust status. Current class: ${tile.className}, New status: ${isTrusted ? 'Trusted' : 'Untrusted'}`);

    // Define the possible border classes
    const trustedBorder = "border-custom-teal";
    const untrustedBorder = "border-custom-yellow";

    // Remove existing border classes
    tile.classList.remove(trustedBorder, untrustedBorder);

    // Apply the appropriate border class and update data attribute
    if (isTrusted) {
        tile.classList.add(trustedBorder);
        tile.dataset.trustStatus = "trusted";
        console.log(`Applied 'trusted' border class to tile.`);
    } else {
        tile.classList.add(untrustedBorder);
        tile.dataset.trustStatus = "untrusted";
        console.log(`Applied 'untrusted' border class to tile.`);
    }

    updateWriteAllButtonState(); // Update button state after operation

}



/**
 * Recomputes and updates a character's tile trust status.
 * @param {number|string} characterID - The ID of the character.
 */
function recomputeAndUpdateTileTrustStatus(characterID) {
    // Convert characterID to number to match the data type in TabulatorIdentities
    const numericCharacterID = Number(characterID);
    if (isNaN(numericCharacterID)) {
        console.warn(`Invalid character ID: ${characterID}`);
        return;
    }

    const character = TabulatorIdentities.find(char => char.CharacterID === numericCharacterID);
    if (character) {
        const tile = document.querySelector(`.character-tile[data-id="${character.CharacterID}"]`);
        if (tile) {
            const isTrusted = isCharacterTrusted(character);
            updateTileTrustStatus(tile, isTrusted);
            console.log(`Updated trust status for CharacterID ${character.CharacterID}: ${isTrusted ? 'Trusted' : 'Untrusted'}`);
        }
    } else {
        console.warn(`Character with ID ${characterID} not found.`);
    }
}

/**
 * Initializes a character tile
 * @param {object} character - The character data object
 * @returns {HTMLElement} - The character tile element
 */
/**
 * Initializes a character tile
 * @param {object} character - The character data object
 * @returns {HTMLElement} - The character tile element
 */
function initializeCharacterTile(character) {
    const tile = document.createElement("div");
    tile.className = "character-tile border-2 rounded-lg p-4 flex flex-col items-center"; // Base Tailwind classes

    tile.dataset.id = character.CharacterID; // Ensure correct property name

    // Determine if the character is trusted
    const isTrusted = character.IsTrusted; // Assuming this is a boolean

    // Use the updateTileTrustStatus function to set border and data attribute
    updateTileTrustStatus(tile, isTrusted);

    const img = document.createElement("img");
    img.src = character.Portrait;
    img.alt = `${character.CharacterName} Portrait`;
    img.className = "character-portrait";

    const name = document.createElement("div");
    name.className = "character-name";
    name.innerText = character.CharacterName;

    const button = document.createElement("button");
    button.className = "write-contacts-btn bg-custom-teal hover:bg-custom-teal-dark text-gray-900 p-2 rounded-full focus:outline-none focus:ring-2 focus:ring-custom-yellow";
    button.title = "Write Contacts";
    button.setAttribute("data-tooltip", "Write Contacts");
    button.setAttribute("aria-label", "Write Contacts");
    button.innerHTML = '<i class="fas fa-pen" aria-hidden="true"></i>';

    button.addEventListener("click", (e) => {
        e.stopPropagation();
        writeContacts(character.CharacterID);
    });

    tile.appendChild(img);
    tile.appendChild(name);
    tile.appendChild(button);
    return tile;
}

function initializeCharacterTiles() {
    const characterContainer = document.getElementById("character-container");
    if (!characterContainer) {
        console.error(`Character container with ID "character-container" not found.`);
        return;
    }
    characterContainer.innerHTML = ""; // Clear existing tiles to prevent duplicates

    // Calculate the appropriate number of columns based on the number of tiles
    const numCharacters = TabulatorIdentities.length;
    const columns = Math.min(numCharacters, 4); // Limit to a max of 4 columns

    // Remove any grid-related classes that might interfere
    characterContainer.className = "";

    // Apply necessary classes for layout
    characterContainer.classList.add("grid", "gap-6", "justify-items-center", "mx-auto");

    // Set grid-template-columns to match the number of character tiles
    characterContainer.style.display = "grid";
    characterContainer.style.gridTemplateColumns = `repeat(${columns}, minmax(150px, 1fr))`;
    characterContainer.style.maxWidth = "fit-content"; // Ensures the container wraps tightly around tiles

    TabulatorIdentities.forEach(character => {
        console.log("Initializing tile for character:", character); // Debugging line
        const tile = initializeCharacterTile(character);
        characterContainer.appendChild(tile);

        // Click event to add trusted characters from untrusted status
        tile.addEventListener("click", () => {
            console.log(`Tile clicked for CharacterID: ${character.CharacterID}, Trust Status: ${tile.dataset.trustStatus}`);
            if (tile.dataset.trustStatus === 'untrusted' && activeRequests === 0) {
                console.log("Adding to trusted list...");
                addEntity('trusted', 'character', character.CharacterID.toString()); // Convert to string
            } else {
                console.log("Click ignored. Either trusted or request in progress.");
            }
        });
    });

    updateWriteAllButtonState();
}




/**
 * Adds an entity based on trustStatus and entityType using a single identifier.
 * @param {string} trustStatus - 'trusted' or 'untrusted'.
 * @param {string} entityType - 'character' or 'corporation'.
 * @param {string|number} identifier - Character/Corporation ID or Name.
 */
function addEntity(trustStatus, entityType, identifier) {
    console.log("Adding entity:", trustStatus, entityType, identifier);

    const serverEndpoint = getServerEndpoint(trustStatus, entityType, 'add');
    if (!serverEndpoint) {
        console.error("Invalid trustStatus or entityType provided to addEntity.");
        toastr.error("An unexpected error occurred.");
        return;
    }

    // Ensure identifier is always a string
    const identifierStr = String(identifier);

    const oppositeTrustStatus = trustStatus === 'trusted' ? 'untrusted' : 'trusted';

    // Check if the entity is already in the opposite list
    const isInOppositeList = trustStatus === 'trusted'
        ? (entityType === 'character' ? isCharacterInUntrustedList(identifier) : isCorporationInUntrustedList(identifier))
        : (entityType === 'character' ? isCharacterInTrustedList(identifier) : isCorporationInTrustedList(identifier));

    if (isInOppositeList) {
        toastr.warning(`${capitalize(entityType)} already exists in the ${oppositeTrustStatus} list.`);
        console.warn(`${capitalize(entityType)} with identifier ${identifierStr} is already in ${oppositeTrustStatus} list.`);
        return; // Prevent adding to the current list
    }

    // Prepare payload with a single identifier field
    const payload = {
        identifier: identifierStr
    };

    showLoading();
    const toggleBtn = document.getElementById("toggle-contacts-btn");
    toggleButtonState(toggleBtn, true);

    fetchWithHandling(serverEndpoint, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload)
    })
        .then(data => {
            console.log(`Entity added:`, data);

            // Update local data and the table
            additionUpdateLocalData(trustStatus, entityType, data);
            const tableId = `${trustStatus}-${entityType}s-table`;
            addRowToTable(tableId, data);

            // Recompute and update the trust status of affected characters
            if (entityType === 'character') {
                recomputeAndUpdateTileTrustStatus(data.CharacterID);
            } else if (entityType === 'corporation') {
                // Update all characters belonging to this corporation
                TabulatorIdentities
                    .filter(char => char.CorporationID === data.CorporationID)
                    .forEach(char => recomputeAndUpdateTileTrustStatus(char.CharacterID));
            }

        })
        .catch(error => {
            console.error(`Error adding entity: ${error}`);
            toastr.error(`Failed to add ${entityType}. ${error.message}`);
        })
        .finally(() => {
            hideLoading();
            toggleButtonState(toggleBtn, false);
        });
}

/**
 * Updates local data by adding the new entity to the appropriate list
 * @param {string} trustStatus - 'trusted' or 'untrusted'
 * @param {string} entityType - 'character' or 'corporation'
 * @param {object} data - The data object of the entity to add
 */
function additionUpdateLocalData(trustStatus, entityType, data) {
    let targetList;
    if (trustStatus === 'trusted' && entityType === 'character') {
        targetList = TrustedCharacters;
    } else if (trustStatus === 'trusted' && entityType === 'corporation') {
        targetList = TrustedCorporations;
    } else if (trustStatus === 'untrusted' && entityType === 'character') {
        targetList = UntrustedCharacters;
    } else if (trustStatus === 'untrusted' && entityType === 'corporation') {
        targetList = UntrustedCorporations;
    } else {
        console.warn(`Unknown trustStatus (${trustStatus}) or entityType (${entityType})`);
        return;
    }

    // Check for duplicates
    const idField = entityType === 'character' ? 'CharacterID' : 'CorporationID';
    const exists = targetList.some(entity => entity[idField] === data[idField]);

    if (!exists) {
        targetList.push(data);
        console.log(`Added ${entityType} with ID ${data[idField]} to ${trustStatus} list.`);
    } else {
        console.warn(`${entityType} with ID ${data[idField]} already exists in ${trustStatus} list.`);
    }
}

/**
 * Adds a row to the specified table
 * @param {string} tableId - The ID of the table
 * @param {object} data - The data object to add as a row
 */
function addRowToTable(tableId, data) {
    const targetTable = tables[tableId];
    if (targetTable) {
        const rowID = data.CharacterID || data.CorporationID;
        const existingRow = targetTable.getRow(rowID);

        if (!existingRow) {
            targetTable.addRow(data)
                .then(() => {
                    const entityType = tableId.includes('character') ? 'Character' : 'Corporation';
                    const trustStatus = tableId.split('-')[0];
                    toastr.success(`Added ${data[`${entityType}Name`]} to ${capitalize(trustStatus)} list.`);

                    // Determine if the table should be displayed based on current view
                    const isTrustedTable = tableId.startsWith('trusted-');
                    const isUntrustedTable = tableId.startsWith('untrusted-');

                    if (isTrustedTable && !isShowingUntrusted) {
                        // Trusted tables should be visible when showing trusted view
                        toggleSectionDisplay(tableId, true);
                    } else if (isUntrustedTable && isShowingUntrusted) {
                        // Untrusted tables should be visible when showing untrusted view
                        toggleSectionDisplay(tableId, true);
                    } else {
                        // Tables not in the current view should remain hidden
                        toggleSectionDisplay(tableId, false);
                    }

                    resizeTabulatorTable(tableId);
                })
                .catch(error => {
                    console.error(`Error adding row to ${tableId}:`, error);
                    const entityType = tableId.includes('character') ? 'Character' : 'Corporation';
                    const trustStatus = tableId.split('-')[0];
                    toastr.error(`Error updating ${capitalize(trustStatus)} ${entityType} table.`);
                });
        } else {
            const entityType = tableId.includes('character') ? 'Character' : 'Corporation';
            const trustStatus = tableId.split('-')[0];
            toastr.warning(`${entityType} already exists in the ${capitalize(trustStatus)} list.`);
        }
    }
}

/**
 * Handle form submission by extracting the identifier and performing add/remove operations.
 * @param {Event} e - The submit event.
 * @param {string} trustStatus - 'trusted' or 'untrusted'.
 * @param {string} entityType - 'character' or 'corporation'.
 */
function handleFormSubmission(e, trustStatus, entityType) {
    e.preventDefault();

    const inputIdMap = {
        'trusted-character': 'trusted-character-identifier',
        'trusted-corporation': 'trusted-corporation-identifier',
        'untrusted-character': 'untrusted-character-identifier',
        'untrusted-corporation': 'untrusted-corporation-identifier',
    };

    const inputId = inputIdMap[`${trustStatus}-${entityType}`];
    const inputElement = document.getElementById(inputId);
    if (!inputElement) {
        toastr.error(`Input element with ID "${inputId}" not found.`);
        console.error(`Input element with ID "${inputId}" not found.`);
        return;
    }

    const identifier = inputElement.value.trim();

    console.log(`Form Submission - Trust Status: ${trustStatus}, Entity Type: ${entityType}, Identifier: "${identifier}"`);

    if (!identifier) {
        toastr.error(`${capitalize(entityType)} needs a name or id.`);
        return;
    }

    // Optional: Additional validation based on expected formats
    if (isNumeric(identifier) && parseInt(identifier, 10) <= 0) {
        toastr.error("Identifier must be a positive number.");
        return;
    }

    // Prevent adding an entity that already exists in the opposite list
    const isInOppositeList = trustStatus === 'trusted'
        ? (entityType === 'character' ? isCharacterInUntrustedList(identifier) : isCorporationInUntrustedList(identifier))
        : (entityType === 'character' ? isCharacterInTrustedList(identifier) : isCorporationInTrustedList(identifier));

    if (isInOppositeList) {
        const statusMessage = trustStatus === 'trusted' ? 'untrusted' : 'trusted';
        toastr.warning(`${identifier} is already in the ${statusMessage} list.`);
        return;
    }

    // Call addEntity with the single identifier
    addEntity(trustStatus, entityType, identifier);
    inputElement.value = '';  // Clear the input field after submission
}

/**
 * Setup event listeners for forms using a centralized handler
 */
function setupFormEventListeners() {
    const forms = [
        { id: "add-trusted-character-form", trustStatus: 'trusted', entityType: 'character' },
        { id: "add-untrusted-character-form", trustStatus: 'untrusted', entityType: 'character' },
        { id: "add-trusted-corporation-form", trustStatus: 'trusted', entityType: 'corporation' },
        { id: "add-untrusted-corporation-form", trustStatus: 'untrusted', entityType: 'corporation' },
    ];

    forms.forEach(({ id, trustStatus, entityType }) => {
        const form = document.getElementById(id);
        if (form) {
            form.addEventListener("submit", (e) => handleFormSubmission(e, trustStatus, entityType));
        } else {
            console.warn(`Form with ID "${id}" not found.`);
        }
    });
}

/**
 * Initializes all Tabulator tables
 */
function initializeAllTabulatorTables() {
    // Define table configurations
    const tableConfigs = [
        {
            tableId: "trusted-characters-table",
            indexField: "CharacterID",
            data: TrustedCharacters,
            columns: [
                {
                    title: "Status",
                    field: "IsOnCouch", // This should match the field name in the data
                    formatter: (cell, formatterParams) => {
                        const isOnCouch = cell.getValue();
                        return isOnCouch
                            ? '<i class="fas fa-couch text-yellow-500" title="On the Couch"></i>'
                            : '<i class="fas fa-user text-green-500" title="Member"></i>';
                    },
                    hozAlign: "center",
                    minWidth: 50,
                    cellClick: function (e, cell) {
                        const rowData = cell.getRow().getData();
                        const currentStatus = rowData.IsOnCouch;

                        // Toggle the status
                        const newStatus = !currentStatus;

                        // Update the backend
                        updateStatus(rowData.CharacterID, newStatus, "character")
                            .then(() => {
                                // Update the cell value and refresh the table
                                rowData.IsOnCouch = newStatus;
                                cell.getRow().update(rowData);
                                toastr.success("Status updated successfully.");
                            })
                            .catch(error => {
                                console.error("Failed to update status:", error);
                                toastr.error("Failed to update status. Please try again.");
                            });
                    }
                },
                { title: "Character Name", field: "CharacterName", headerSort: true, minWidth: 100 },
                { title: "Corporation", field: "CorporationName", headerSort: true, minWidth: 100 },
                {
                    title: "Comment",
                    field: "Comment",
                    editor: "input",
                    editable: true,
                    minWidth: 150
                },
                {
                    title: "Remove",
                    formatter: "buttonCross",
                    hozAlign: "center",
                    headerSort: false,
                    minWidth: 50,
                    cellClick: function (e, cell) {
                        const rowData = cell.getRow().getData();
                        const characterID = rowData.CharacterID;
                        const characterName = rowData.CharacterName;
                        console.log(`Removing trusted character with ID: ${characterID}, Name: ${characterName}`);

                        Swal.fire({
                            title: `Remove Character?`,
                            text: `Do you want to stop trusting "${characterName}"?`,
                            icon: 'warning',
                            showCancelButton: true,
                            confirmButtonText: 'Yes',
                            cancelButtonText: 'No',
                        }).then((result) => {
                            if (result.isConfirmed) {
                                removeEntity('trusted', 'character', characterID.toString());
                            }
                        });
                    }
                }
            ]
        },
        {
            tableId: "trusted-corporations-table",
            indexField: "CorporationID",
            data: TrustedCorporations,
            columns: [
                {
                    title: "Status",
                    field: "IsOnCouch", // This should match the field name in the data
                    formatter: (cell, formatterParams) => {
                        const isOnCouch = cell.getValue();
                        return isOnCouch
                            ? '<i class="fas fa-couch text-yellow-500" title="On the Couch"></i>'
                            : '<i class="fas fa-user text-green-500" title="Member"></i>';
                    },
                    hozAlign: "center",
                    minWidth: 50,
                    cellClick: function (e, cell) {
                        const rowData = cell.getRow().getData();
                        const currentStatus = rowData.IsOnCouch;

                        // Toggle the status
                        const newStatus = !currentStatus;

                        // Update the backend
                        updateStatus(rowData.CharacterID, newStatus, "character")
                            .then(() => {
                                // Update the cell value and refresh the table
                                rowData.IsOnCouch = newStatus;
                                cell.getRow().update(rowData);
                                toastr.success("Status updated successfully.");
                            })
                            .catch(error => {
                                console.error("Failed to update status:", error);
                                toastr.error("Failed to update status. Please try again.");
                            });
                    }
                },
                { title: "Corporation Name", field: "CorporationName", headerSort: true, minWidth: 100 },
                { title: "Alliance Name", field: "AllianceName", headerSort: true, minWidth: 100 },
                {
                    title: "Comment",
                    field: "Comment",
                    editor: "input",
                    editable: true,
                    minWidth: 150
                },
                {
                    title: "Remove",
                    formatter: "buttonCross",
                    hozAlign: "center",
                    headerSort: false,
                    minWidth: 50,
                    cellClick: function (e, cell) {
                        const rowData = cell.getRow().getData();
                        const corporationID = rowData.CorporationID;
                        const corporationName = rowData.CorporationName;
                        console.log(`Removing trusted corporation with ID: ${corporationID}, Name: ${corporationName}`);

                        Swal.fire({
                            title: `Remove Corporation?`,
                            text: `Do you want to stop trusting "${corporationName}"?`,
                            icon: 'warning',
                            showCancelButton: true,
                            confirmButtonText: 'Yes',
                            cancelButtonText: 'No',
                        }).then((result) => {
                            if (result.isConfirmed) {
                                removeEntity('trusted', 'corporation', corporationID.toString());
                            }
                        });
                    }
                }
            ]
        },
        {
            tableId: "untrusted-characters-table",
            indexField: "CharacterID",
            data: UntrustedCharacters,
            columns: [
                { title: "Character Name", field: "CharacterName", headerSort: true, minWidth: 250 },
                { title: "Added By", field: "AddedBy", headerSort: true, minWidth: 250 },
                { title: "Corporation", field: "CorporationName", headerSort: true, minWidth: 250 },
                {
                    title: "Comment",
                    field: "Comment",
                    editor: "input", // Makes the cell editable
                    editable: true, // Ensures it's editable by the user
                    minWidth: 300,
                },
                {
                    title: "Remove",
                    formatter: "buttonCross",
                    hozAlign: "center",
                    headerSort: false,
                    minWidth: 50,
                    cellClick: function (e, cell) {
                        const rowData = cell.getRow().getData();
                        const characterName = rowData.CharacterName;
                        const characterID = rowData.CharacterID;
                        console.log(`Removing untrusted character with ID: ${characterID}, Name: ${characterName}`);

                        // Use SweetAlert2 for Confirmation
                        Swal.fire({
                            title: `Remove Character?`,
                            text: `Has everyone updated their standings for "${characterName}"?`,
                            icon: 'warning',
                            showCancelButton: true,
                            confirmButtonText: 'Yes',
                            cancelButtonText: 'No',
                            customClass: {
                                popup: 'bg-gray-800 text-gray-200', // Tailwind classes for modal background and text
                                title: 'font-semibold text-xl', // Tailwind classes for title styling
                                content: 'text-gray-300', // Tailwind classes for content/body text
                                confirmButton: 'bg-teal-500 hover:bg-teal-600 text-gray-900 p-2 rounded-full focus:outline-none focus:ring-2 focus:ring-yellow-400',
                                cancelButton: 'bg-red-500 hover:bg-red-600 text-gray-900 p-2 rounded-full focus:outline-none focus:ring-2 focus:ring-yellow-400',
                                actions: 'flex justify-center space-x-4', // Tailwind classes for button container
                                // Optionally, customize the icon
                                icon: 'text-yellow-400' // Tailwind class for icon color
                            }
                        }).then((result) => {
                            if (result.isConfirmed) {
                                removeEntity('untrusted', 'character', characterID.toString()); // Convert to string
                            }
                        });
                    }
                }
            ]
        },
        {
            tableId: "untrusted-corporations-table",
            indexField: "CorporationID",
            data: UntrustedCorporations,
            columns: [
                { title: "Corporation Name", field: "CorporationName", headerSort: true, minWidth: 250 },
                { title: "Added By", field: "AddedBy", headerSort: true, minWidth: 250 },
                { title: "Alliance Name", field: "AllianceName", headerSort: true, minWidth: 250 },
                {
                    title: "Comment",
                    field: "Comment",
                    editor: "input", // Makes the cell editable
                    editable: true, // Ensures it's editable by the user
                    minWidth: 300,
                },
                {
                    title: "Remove",
                    formatter: "buttonCross",
                    hozAlign: "center",
                    headerSort: false,
                    minWidth: 50,
                    cellClick: function (e, cell) {
                        const rowData = cell.getRow().getData();
                        const corporationID = rowData.CorporationID;
                        const corporationName = rowData.CorporationName;
                        console.log(`Removing untrusted corporation with ID: ${corporationID}, Name: ${corporationName}`);

                        // Use SweetAlert2 for Confirmation
                        Swal.fire({
                            title: `Remove Corporation?`,
                            text: `Has everyone updated their standings for "${corporationName}"?`,
                            icon: 'warning',
                            showCancelButton: true,
                            confirmButtonText: 'Yes',
                            cancelButtonText: 'No',
                            customClass: {
                                popup: 'bg-gray-800 text-gray-200', // Tailwind classes for modal background and text
                                title: 'font-semibold text-xl', // Tailwind classes for title styling
                                content: 'text-gray-300', // Tailwind classes for content/body text
                                confirmButton: 'bg-teal-500 hover:bg-teal-600 text-gray-900 p-2 rounded-full focus:outline-none focus:ring-2 focus:ring-yellow-400',
                                cancelButton: 'bg-red-500 hover:bg-red-600 text-gray-900 p-2 rounded-full focus:outline-none focus:ring-2 focus:ring-yellow-400',
                                actions: 'flex justify-center space-x-4', // Tailwind classes for button container
                                // Optionally, customize the icon
                                icon: 'text-yellow-400' // Tailwind class for icon color
                            }
                            }).then((result) => {
                            if (result.isConfirmed) {
                                removeEntity('untrusted', 'corporation', corporationID.toString()); // Convert to string
                            }
                        });
                    },
                }
            ]
        }
    ];

    // Initialize each table using the generic function
    tableConfigs.forEach(config => {
        initializeTabulatorTable(config.tableId, config.indexField, config.data, config.columns);
    });

    // Setup Mutation Observers for all tables
    const allTableIds = tableConfigs.map(config => config.tableId);
    setupMutationObservers(allTableIds);
}

async function updateStatus(entityID, isOnCouch, entityType) {
    const url = `/update-is-on-couch`;

    const payload = {
        id: entityID,
        isOnCouch,
        tableId: entityType === "character" ? "trusted-characters-table" : "trusted-corporations-table"
    };

    try {
        const response = await fetch(url, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload)
        });

        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.error || "Failed to update status.");
        }

        return response.json();
    } catch (error) {
        console.error(`Error updating status:`, error);
        throw error;
    }
}

/**
 * Shows the untrusted tables and redraws them
 */
function showUntrustedTables() {
    const untrustedSections = [
        "untrusted-characters-table",
        "untrusted-corporations-table",
        "add-untrusted-character-section",
        "add-untrusted-corporation-section"
    ];

    untrustedSections.forEach(sectionId => {
        const element = document.getElementById(sectionId);
        if (element) {
            element.style.visibility = "visible";
            element.style.opacity = "1";
            element.style.position = "relative";
        }
    });

    // Force redraw for untrusted tables after making them visible
    if (tables["untrusted-characters-table"]) {
        tables["untrusted-characters-table"].redraw(true);
    }
    if (tables["untrusted-corporations-table"]) {
        tables["untrusted-corporations-table"].redraw(true);
    }
}

/**
 * Hides the untrusted tables by setting visibility to hidden
 */
function hideUntrustedTables() {
    const untrustedSections = [
        "untrusted-characters-table",
        "untrusted-corporations-table",
        "add-untrusted-character-section",
        "add-untrusted-corporation-section"
    ];

    untrustedSections.forEach(sectionId => {
        const element = document.getElementById(sectionId);
        if (element) {
            element.style.visibility = "hidden";
            element.style.opacity = "0";
            element.style.position = "absolute";
        }
    });
}



/**
 * Sends the updated comment to the backend.
 * @param {number} id - ID to be commented on
 * @param {string} comment - The updated comment text
 * @param {string} tableId - ID of the table the comment is in
 */
async function updateComment(id, comment, tableId) {
    const url = `/update-comment`; // Replace with your actual endpoint URL

    try {
        showLoading();
        const response = await fetch(url, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ id, comment, tableId })
        });

        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.error || "Failed to save comment.");
        }

        toastr.success("Comment saved successfully.");
    } catch (error) {
        console.error(`Error saving comment: ${error}`);
        toastr.error("Failed to save comment. " + error.message);
    } finally {
        hideLoading();
    }
}

/**
 * Removes an entity based on trustStatus and entityType using a single identifier.
 * @param {string} trustStatus - 'trusted' or 'untrusted'.
 * @param {string} entityType - 'character' or 'corporation'.
 * @param {string|number} identifier - Character/Corporation ID or Name.
 */
async function removeEntity(trustStatus, entityType, identifier) {
    console.log("Removing entity:", trustStatus, entityType, identifier);

    const serverEndpoint = getServerEndpoint(trustStatus, entityType, 'remove');
    if (!serverEndpoint) {
        console.error("Invalid trustStatus or entityType provided to removeEntity.");
        toastr.error("An unexpected error occurred.");
        return;
    }

    // Ensure identifier is always a string
    const identifierStr = String(identifier);

    // Prepare payload
    const payload = { identifier: identifierStr };
    console.log("Removing entity with payload:", payload, "Endpoint:", serverEndpoint);

    try {
        showLoading(); // Show the loading indicator

        // Perform the fetch request
        const response = await fetchWithHandling(serverEndpoint, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload)
        });

        console.log(`Entity removed successfully: ${identifierStr}`);

        const tableId = `${trustStatus}-${entityType}s-table`;
        removeRowFromTable(tableId, identifierStr);

        // Update local data and the table
        removeUpdateLocalData(trustStatus, entityType, identifierStr);

        if (trustStatus === 'trusted') {
            addEntity("untrusted", entityType, identifier);
        }

        // toastr.success(`Successfully removed ${capitalize(entityType)}.`);
    } catch (error) {
        console.error(`Error removing ${entityType}:`, error);
        toastr.error(`Failed to remove ${capitalize(entityType)}. ${error.message}`);
    } finally {
        hideLoading(); // Hide the loading indicator
    }
}

/**
 * Removes entity data from the appropriate local list.
 * @param {string} trustStatus - 'trusted' or 'untrusted'
 * @param {string} entityType - 'character' or 'corporation'
 * @param {string|number} identifier - The identifier of the entity to remove
 */
function removeUpdateLocalData(trustStatus, entityType, identifier) {
    const identifierNum = Number(identifier); // Ensures compatibility if identifier is a string

    // Determine which list to modify directly
    if (trustStatus === 'trusted' && entityType === 'character') {
        TrustedCharacters = TrustedCharacters.filter(entity => entity.CharacterID !== identifierNum);
    } else if (trustStatus === 'trusted' && entityType === 'corporation') {
        TrustedCorporations = TrustedCorporations.filter(entity => entity.CorporationID !== identifierNum);
    } else if (trustStatus === 'untrusted' && entityType === 'character') {
        UntrustedCharacters = UntrustedCharacters.filter(entity => entity.CharacterID !== identifierNum);
    } else if (trustStatus === 'untrusted' && entityType === 'corporation') {
        UntrustedCorporations = UntrustedCorporations.filter(entity => entity.CorporationID !== identifierNum);
    } else {
        console.warn(`Unknown trustStatus (${trustStatus}) or entityType (${entityType})`);
        return;
    }

    console.log(`Removed ${entityType} with ID ${identifier} from ${trustStatus} list.`);
}

/**
 * Removes a row from the specified table by identifier.
 * @param {string} tableId - The ID of the table.
 * @param {string|number} identifier - The identifier of the row to remove.
 */
function removeRowFromTable(tableId, identifier) {
    const targetTable = tables[tableId];
    if (targetTable) {
        const numericIdentifier = Number(identifier);
        const entityType = tableId.includes('character') ? 'Character' : 'Corporation';
        const trustStatus = tableId.split('-')[0];
        const data = targetTable.getData().find(row => row.CharacterID === numericIdentifier || row.CorporationID === numericIdentifier);

        if (!data) {
            console.warn(`No data found for identifier ${identifier} in table ${tableId}.`);
            toastr.warning(`No matching ${entityType} found to remove.`);
            return;
        }

        targetTable.deleteRow(numericIdentifier)
            .then(() => {
                resizeTabulatorTable(tableId);
                if (targetTable.getData().length === 0) {
                    toggleSectionDisplay(tableId, false);
                }
                toastr.success(`Removed ${data[`${entityType}Name`]} from ${capitalize(trustStatus)} list.`);
            })
            .catch(error => {
                console.error(`Error deleting row from ${tableId}:`, error);
                const entityType = tableId.includes('character') ? 'Character' : 'Corporation';
                const trustStatus = tableId.split('-')[0];
                toastr.error(`Error updating ${data[`${entityType}Name`]} in ${capitalize(trustStatus)} ${entityType} table.`);
            });
    } else {
        console.warn(`Table with ID "${tableId}" not found.`);
    }
}

/**
 * Capitalizes the first letter of a string
 * @param {string} str - The string to capitalize.
 * @returns {string} - The capitalized string.
 */
function capitalize(str) {
    if (typeof str !== 'string') return '';
    return str.charAt(0).toUpperCase() + str.slice(1);
}

/**
 * Checks if a string is numeric.
 * @param {string} str - The string to check.
 * @returns {boolean} - True if numeric, else false.
 */
function isNumeric(str) {
    return /^\d+$/.test(str);
}

/**
 * Setup Mutation Observers for tables
 * @param {Array} tableIds - Array of table IDs to observe
 */
function setupMutationObservers(tableIds) {
    tableIds.forEach(tableId => {
        const tableElement = document.getElementById(tableId);
        if (!tableElement) {
            console.warn(`Table element with ID "${tableId}" not found for MutationObserver.`);
            return;
        }

        new MutationObserver(function(mutations) {
            mutations.forEach(function(mutation) {
                if (mutation.attributeName === "style") {
                    const tableDisplay = window.getComputedStyle(tableElement).display;
                    console.log(`Style mutation detected for ${tableId}: ${tableDisplay}`);

                    const isTrustedTable = tableId.startsWith('trusted-');
                    const isUntrustedTable = tableId.startsWith('untrusted-');

                    if (isTrustedTable && !isShowingUntrusted) {
                        // Show trusted tables
                        toggleSectionDisplay(tableId, true);
                    } else if (isUntrustedTable && isShowingUntrusted) {
                        // Show untrusted tables
                        toggleSectionDisplay(tableId, true);
                    } else {
                        // Hide tables not currently in view
                        toggleSectionDisplay(tableId, false);
                    }
                }
            });
        }).observe(tableElement, { attributes: true });
    });
}

/**
 * Initializes a Tabulator table with given parameters
 * @param {string} tableId - The ID of the table container
 * @param {string} indexField - The field to use as row index
 * @param {Array} data - The initial data array
 * @param {Array} columns - The column definitions
 */
function initializeTabulatorTable(tableId, indexField, data, columns) {
    tables[tableId] = new Tabulator(`#${tableId}`, {
        index: indexField,
        data: data || [],
        layout: "fitData",
        responsiveLayout: "hide",
        placeholder: `No ${capitalize(tableId.split('-')[1].slice(0, -1))}s`,
        columns: columns,
        rowAdded: function () {
            resizeTabulatorTable(tableId);
        },
        rowDeleted: function () {
            resizeTabulatorTable(tableId);
        },
        tableBuilt: function () {
            // Ensure resize happens after the table is fully built
            resizeTabulatorTable(tableId);
        },
        cellEdited: function (cell) {
            // Ensure this is for the "Comment" field
            if (cell.getColumn().getField() === "Comment") {
                const rowData = cell.getRow().getData();
                const updatedComment = cell.getValue();

                // Determine if this is a character or corporation table
                const isCharacterTable = tableId.includes("character");
                const entityId = isCharacterTable ? rowData.CharacterID : rowData.CorporationID;

                // Log for debugging (optional)
                console.log(`Updating comment for ${isCharacterTable ? "Character" : "Corporation"} ID: ${entityId}, Comment: ${updatedComment}`);
                console.log(`Table ID: ${tableId}`);
                // Call backend function to update the comment
                updateComment(entityId, updatedComment, tableId);
            }
        }
    });
}

/**
 * Toggle Button Event Listener
 * Switches between showing trusted and untrusted sections
 */
function setupToggleButton() {
    const toggleBtn = document.getElementById("toggle-contacts-btn");
    if (toggleBtn) {
        toggleBtn.addEventListener("click", function () {
            console.log("Toggle button clicked.");

            const trustedSections = [
                "trusted-characters-table",
                "trusted-corporations-table",
                "add-trusted-character-section",
                "add-trusted-corporation-section"
            ];

            const untrustedSections = [
                "untrusted-characters-table",
                "untrusted-corporations-table",
                "add-untrusted-character-section",
                "add-untrusted-corporation-section"
            ];

            const icon = this.querySelector("i");

            if (isShowingUntrusted) {
                console.log("Switching to trusted view...");

                // Show Trusted Tables and Forms
                toggleMultipleSections(trustedSections, true);
                hideUntrustedTables();

                icon.classList.remove("fa-toggle-off");
                icon.classList.add("fa-toggle-on");
                this.title = "Show Untrusted Contacts";

                isShowingUntrusted = false;

            } else {
                console.log("Switching to untrusted view...");

                // Show Untrusted Tables and Forms
                toggleMultipleSections(trustedSections, false);
                showUntrustedTables();

                icon.classList.remove("fa-toggle-on");
                icon.classList.add("fa-toggle-off");
                this.title = "Show Trusted Contacts";

                isShowingUntrusted = true;
            }

            // Force redraw of all visible tables after toggling
            setTimeout(() => {
                const tableIds = [
                    "trusted-characters-table",
                    "trusted-corporations-table",
                    "untrusted-characters-table",
                    "untrusted-corporations-table"
                ];

                tableIds.forEach(tableId => {
                    const tableElement = document.getElementById(tableId);
                    if (tableElement && window.getComputedStyle(tableElement).visibility === "visible") {
                        tables[tableId].redraw(true);  // Force redraw for visible tables
                    }
                });
            }, 200); // Delay to allow the DOM to update
        });
    } else {
        console.error(`Toggle button with ID "toggle-contacts-btn" not found.`);
    }
}


/**
 * Initialize Everything After DOM is Loaded
 */
document.addEventListener('DOMContentLoaded', function () {
    // Set initial state of the view
    isShowingUntrusted = false;  // Ensure trusted tables visible and untrusted hidden.
    hideUntrustedTables();

    // Hide untrusted tables and sections immediately on page load
    const untrustedSections = [
        "untrusted-characters-table",
        "untrusted-corporations-table",
        "add-untrusted-character-section",
        "add-untrusted-corporation-section"
    ];
    toggleMultipleSections(untrustedSections, false);

    // Initialize all Tabulator tables
    initializeAllTabulatorTables();

    // Show trusted sections if they have data
    const trustedSections = [
        "trusted-characters-table",
        "trusted-corporations-table"
    ];
    trustedSections.forEach(tableId => {
        const hasData = tables[tableId].getData().length > 0;
        toggleSectionDisplay(tableId, hasData);
    });

    // Initialize character tiles
    initializeCharacterTiles();

    // Setup Toggle Button Event Listener
    setupToggleButton();

    // Initial resizing of tables on page load
    setTimeout(() => {
        const initialTables = [
            "trusted-characters-table",
            "trusted-corporations-table"
            // Untrusted tables are hidden on page load
        ];

        initialTables.forEach(tableId => {
            const tableElement = document.getElementById(tableId);
            if (tableElement && window.getComputedStyle(tableElement).display !== "none") {
                resizeTabulatorTable(tableId);
            }
        });

        // Untrusted tables and sections are already hidden earlier
    }, 200);

    setupWriteAllButton();

    // Setup all form event listeners
    setupFormEventListeners();
    hideLoading();
});
