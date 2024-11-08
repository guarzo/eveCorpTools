function initSaveSplit() {
    const container = document.getElementById("battle-report-container");
    container.innerHTML = `
        <div class="input-container text-center">
            <textarea id="battleReport" 
                      class="w-full h-12 p-4 text-gray-200 bg-gray-700 rounded-lg resize-none focus:outline-none overflow-y-hidden"
                      placeholder="Enter battle report..." rows="1"></textarea>
        </div>
        <div class="button-container mt-4 text-center">
            <button id="saveSplitButton" class="px-4 py-2 bg-teal-500 text-gray-100 font-semibold rounded shadow hover:bg-teal-400">
                Save Split
            </button>
        </div>`;



    document.getElementById("saveSplitButton").addEventListener("click", saveLootSplit);
}

async function saveLootSplit() {
    const totalBuyPriceElement = document.getElementById("lootInput");
    const totalBuyPrice = parseFloat(totalBuyPriceElement.innerText.replace(/,/g, '').replace(' ISK', '')); // Remove commas and convert to number
    const splitDetails = {};
    const battleReportElement = document.getElementById("battleReport");
    const battleReport = battleReportElement ? battleReportElement.value : '';

    if (totalBuyPrice <= 0 || isNaN(totalBuyPrice)) {
        alert("Total buy price must be greater than 0.");
        return;
    }

    if (!battleReport) {
        alert("Battle Report is required.");
        return;
    }

    // Extract values for Scout, Involved pilots, and Corp from valuesContainer
    document.getElementById("valuesContainer").querySelectorAll('p').forEach(p => {
        const textContent = p.textContent.trim();

        // Match Scout and Involved entries with names and amounts
        const match = textContent.match(/(Scout|Involved):\s*([^-\d]+)?\s*-\s*([\d,]+)/);

        if (match) {
            const role = match[1];
            const name = match[2]?.trim() || ''; // Handle name if it exists
            const amount = match[3].replace(/,/g, ''); // Remove commas from amount

            if (name) {
                // For Involved and Scout, store with the name as part of the key
                splitDetails[`${role}: ${name}`] = amount;
            }
        }

        // Separate check for the Corp entry
        if (textContent.startsWith("Corp:")) {
            const corpAmount = textContent.split(":")[1].trim().replace(/,/g, ''); // Extract amount and remove commas
            splitDetails["Corp"] = corpAmount;
        }
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

