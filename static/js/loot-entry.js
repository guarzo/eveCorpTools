function initLootEntry() {
    const container = document.getElementById("loot-entry-container");
    container.innerHTML = `
        <div class="input-container text-center">
            <textarea id="lootEntry" class="w-full p-4 text-gray-200 bg-gray-700 rounded-lg resize-none focus:outline-none" rows="4"
                      placeholder="Enter loot information..." oninput="autoExpand(this)"></textarea>
        </div>
        <div class="button-container mt-4 flex justify-center">
            <button id="appraiseLootButton" class="px-6 py-2 bg-gray-800 text-teal-500 font-semibold rounded-full shadow-lg hover:bg-gray-700 transition duration-300" title="Appraise Loot">
                <i class="fas fa-magic text-teal-500 text-2xl hover:text-teal-300"></i>
            </button>
        </div>`;

    document.getElementById("appraiseLootButton").addEventListener("click", fetchLootPrice);

    const jitaPriceContainer = document.getElementById("jita-price-container");
    jitaPriceContainer.innerHTML = `
        <div class="input-container text-center mt-6">
            <span id="lootInput" class="text-2xl font-semibold text-teal-200">0 ISK</span>
        </div>`;

}

function autoExpand(field) {
    field.style.height = 'inherit';
    const computed = window.getComputedStyle(field);
    field.style.height = (parseInt(computed.getPropertyValue('border-top-width')) +
        parseInt(computed.getPropertyValue('padding-top')) +
        field.scrollHeight +
        parseInt(computed.getPropertyValue('padding-bottom')) +
        parseInt(computed.getPropertyValue('border-bottom-width'))) + 'px';
}

async function fetchLootPrice() {
    const lootEntry = document.getElementById("lootEntry").value;

    if (!lootEntry.trim()) {
        document.getElementById('validationMessage').classList.remove('hidden');
        return;
    } else {
        document.getElementById('validationMessage').classList.add('hidden');
    }

    try {
        const response = await fetch('/appraise-loot', {
            method: 'POST',
            headers: {
                'Content-Type': 'text/plain'
            },
            body: lootEntry
        });

        if (!response.ok) {
            const errorText = await response.text();
            console.error("Error response:", errorText);
            throw new Error("Network response was not ok");
        }


        const data = await response.json();
        const lootInput = document.getElementById("lootInput");
        lootInput.innerText = formatNumber(data.totalBuyPrice) + " ISK";
        lootInput.classList.add('highlight');
        setTimeout(() => lootInput.classList.remove('highlight'), 2000);

        // Show the hidden containers
        document.getElementById('jita-price-container').classList.remove('hidden');
        document.getElementById('pilot-count-container').classList.remove('hidden');
        document.getElementById('scanner-count-container').classList.remove('hidden');
        document.getElementById('calculation-result-container').classList.remove('hidden');
        document.getElementById('battle-report-container').classList.remove('hidden');
        document.getElementById('first-divider').classList.remove('hidden');
        document.getElementById('second-divider').classList.remove('hidden');

        // Initialize calculation result container
        initCalculationResult();

        // Initialize dropdowns before calculations
        populatePilotDropdown();
        initializeScannerDropdown();

        // Add event listeners for dropdown changes
        document.getElementById("pilotCount").addEventListener("change", recalculateSplit);
        document.getElementById("scannerCount").addEventListener("change", recalculateSplit);

        // Initialize battle report container
        initSaveSplit();

        // Calculate values
        calculateValues(data.totalBuyPrice);

        // Scroll to the next section
        document.getElementById('pilot-count-container').scrollIntoView({ behavior: 'smooth' });

        // Now that everything succeeded, show the success message
        toastr.success('Loot appraised successfully!');
        setTimeout(() => {
            if (document.getElementById('valuesContainer')) {
                calculateValues(data.totalBuyPrice);
            } else {
                console.error('valuesContainer not found');
            }
        }, 0); // Ensure DOM is updated before calling

    } catch (error) {
        console.error("Error fetching loot price:", error);
        toastr.error('Failed to appraise loot. Please try again.');
    }
}

function formatNumber(num) {
    return num.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",");
}

document.addEventListener("DOMContentLoaded", function () {
    initLootEntry()
});
