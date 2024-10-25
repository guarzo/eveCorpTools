function initSaveSplit() {
    const container = document.getElementById("battle-report-container");
    container.innerHTML = `
        <div class="input-container">
            <textarea id="battleReport" rows="1" placeholder="Enter battle report..." style="font-size: 1.2em;"></textarea>
        </div>
        <div class="button-container">
            <button id="saveSplitButton">Save Split</button>
        </div>`;

    document.getElementById("saveSplitButton").addEventListener("click", saveLootSplit);
}

async function saveLootSplit() {
    const totalBuyPriceElement = document.getElementById("lootInput");
    const totalBuyPrice = parseFloat(totalBuyPriceElement.innerText.replace(/,/g, '').replace(' ISK', '')); // Remove commas and convert to number
    const valuesContainer = document.getElementById("valuesContainer");
    const splitDetails = {};
    const battleReportElement = document.getElementById("battleReport");
    const battleReport = battleReportElement ? battleReportElement.value : '';

    if (totalBuyPrice <= 0 || isNaN(totalBuyPrice)) {
        alert("Total buy price must be greater than 0.");
        return;
    }

    const scoutBox = document.getElementById('scout').querySelector('.draggable-box');
    const involvedBox = document.getElementById('involved').querySelectorAll('.draggable-box');
    if (!scoutBox && involvedBox.length === 0) {
        alert("At least one pilot must be selected.");
        return;
    }

    if (!battleReport) {
        alert("Battle Report is required.");
        return;
    }

    valuesContainer.querySelectorAll('p').forEach(p => {
        const key = p.textContent.split(':')[0].trim();
        const value = p.textContent.split(':')[1].replace(/[^\d]/g, '').trim();
        splitDetails[key] = value;
    });

    const lootSplit = {
        totalBuyPrice: totalBuyPrice.toString(),
        splitDetails: splitDetails,
        battleReport: battleReport,
        date: new Date().toISOString()
    };

    const jsonString = JSON.stringify(lootSplit);

    try {
        const response = await fetch('/save-loot-split', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: jsonString
        });

        if (response.ok) {
            window.location.href = '/loot-summary';
        } else {
            console.error("Error saving loot split. Server responded with status:", response.status, response.statusText);
        }
    } catch (error) {
        console.error("Error saving loot split:", error);
    }
}
