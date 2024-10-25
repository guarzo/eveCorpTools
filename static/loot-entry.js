function initLootEntry() {
    const container = document.getElementById("loot-entry-container");
    container.innerHTML = `
        <div class="input-container">
            <textarea id="lootEntry" rows="4" placeholder="Enter loot information..." oninput="autoExpand(this)"></textarea>
            <div class="button-container">
                <button id="appraiseLootButton" class="appraise-button">Appraise Loot</button>
            </div>
        </div>`;

    document.getElementById("appraiseLootButton").addEventListener("click", fetchLootPrice);

    const jitaPriceContainer = document.getElementById("jita-price-container");
    jitaPriceContainer.innerHTML = `
        <div class="input-container">
            <span id="lootInput" class="highlight">0 ISK</span>
        </div>`;
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

function autoExpand(field) {
    field.style.height = 'inherit';
    const computed = window.getComputedStyle(field);
    const height = parseInt(computed.getPropertyValue('border-top-width'), 10)
        + parseInt(computed.getPropertyValue('padding-top'), 10)
        + field.scrollHeight
        + parseInt(computed.getPropertyValue('padding-bottom'), 10)
        + parseInt(computed.getPropertyValue('border-bottom-width'), 10);
    field.style.height = height + 'px';
}
