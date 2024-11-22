function populatePilotDropdown() {
    const pilotCountDropdown = document.getElementById("pilotCount");

    // Populate Total Pilots dropdown
    const maxPilots = 20; // Adjust the maximum number of pilots if needed
    pilotCountDropdown.innerHTML = Array.from(
        { length: maxPilots },
        (_, i) => `<option value="${i + 1}">${i + 1}</option>`
    ).join("");

    // Set default value if not already set
    if (!pilotCountDropdown.value) {
        pilotCountDropdown.value = "1";
    }
}

function initializeScannerDropdown() {
    const scannerCountDropdown = document.getElementById("scannerCount");

    // Populate Scanners dropdown with fixed max of 2
    scannerCountDropdown.innerHTML = Array.from(
        { length: 3 }, // Fixed max count of 2 scanners (0, 1, 2)
        (_, i) => `<option value="${i}">${i}</option>`
    ).join("");

    // Set default value if not already set
    scannerCountDropdown.value = scannerCountDropdown.value || "0";

}

document.addEventListener("DOMContentLoaded", function () {
    // Populate both dropdowns on page load
    populatePilotDropdown();
    initializeScannerDropdown();

    // Add event listeners for dropdown changes
    document.getElementById("pilotCount").addEventListener("change", recalculateSplit);
    document.getElementById("scannerCount").addEventListener("change", recalculateSplit);

});
