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
    const totalBuyPrice = parseFloat(totalBuyPriceElement.innerText.replace(/,/g, "").replace(" ISK", ""));
    const totalPilots = parseInt(document.getElementById("pilotCount").value) || 0;
    const Count = parseInt(document.getElementById("Count").value) || 0;
    const battleReport = document.getElementById("battleReport").value;

    if (totalBuyPrice <= 0 || isNaN(totalBuyPrice)) {
        alert("Total buy price must be greater than 0.");
        return;
    }

    if (!battleReport) {
        alert("Battle Report is required.");
        return;
    }

    // Call shared calculation function
    const results = calculateSplit(totalBuyPrice, totalPilots, Count);

    const lootSplit = {
        totalBuyPrice: totalBuyPrice.toString(),
        pilotCount: totalPilots,
        Count,
        involvedCount: totalPilots, // All pilots are involved
        splitDetails: {
            "": results.Payout,
            "Involved": results.pilotPayout,
            "Corporation": results.corpShare,
        },
        battleReport,
        date: new Date().toISOString(),
    };

    try {
        const response = await fetch("/save-loot-split", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(lootSplit),
        });

        if (response.ok) {
            window.location.href = "/loot-summary";
        } else {
            console.error("Error saving loot split:", response.status, response.statusText);
        }
    } catch (error) {
        console.error("Error saving loot split:", error);
    }
}


