function initLootEntry() {
    const container = document.getElementById("loot-entry-container");
    container.innerHTML = `
    <div class="input-container text-center">
        <textarea id="lootEntry" class="w-full p-4 text-gray-200 bg-gray-700 rounded-lg resize-none focus:outline-none" rows="4"
                  placeholder="Enter loot information..." oninput="autoExpand(this)"></textarea>
        <div class="button-container mt-4">
            <button id="appraiseLootButton" class="px-4 py-2 bg-teal-500 text-gray-100 font-semibold rounded shadow hover:bg-teal-400">Appraise Loot</button>
        </div>
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
        setTimeout(() => {
            if (document.getElementById('valuesContainer')) {
                calculateValues(data.totalBuyPrice);
            } else {
                console.error('valuesContainer not found');
            }
        }, 0); // Ensure DOM is updated before calling
    } catch (error) {
        console.error("Error fetching loot price:", error);
    }
}

function formatNumber(num) {
    return num.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",");
}
