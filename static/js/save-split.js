function initSaveSplit() {
    const container = document.getElementById("battle-report-container");
    container.innerHTML = `
        <div class="flex flex-col items-center mt-6">
            <!-- Label for Battle Report -->
            <label for="battleReport" class="text-teal-200 text-lg font-semibold mb-2">
                Battle Report
                <i class="fas fa-info-circle text-gray-400 ml-1" title="Enter a link or brief description"></i>
            </label>

            <!-- Battle Report Input -->
            <input type="text" id="battleReport"
                   class="w-full md:w-1/2 px-4 py-2 text-gray-200 bg-gray-700 rounded-lg focus:outline-none focus:ring-2 focus:ring-teal-500"
                   placeholder="Enter battle report...">

            <!-- Save Split Button -->
            <button id="saveSplitButton" class="mt-4 px-6 py-2 bg-gray-800 text-teal-500 font-semibold rounded-full shadow-lg hover:bg-gray-700 transition duration-300" title="Save Split">
                <i class="fas fa-save text-teal-500 text-2xl"></i>
            </button>
        </div>`;

    document.getElementById("saveSplitButton").addEventListener("click", saveLootSplit);
}


async function saveLootSplit() {
    const totalBuyPriceElement = document.getElementById("lootInput");
    const totalBuyPrice = parseFloat(totalBuyPriceElement.innerText.replace(/,/g, "").replace(" ISK", ""));
    const totalPilots = parseInt(document.getElementById("pilotCount").value) || 0;
    const scannerCount = parseInt(document.getElementById("scannerCount").value) || 0;
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
    const results = calculateSplit(totalBuyPrice, totalPilots, scannerCount);

    const lootSplit = {
        totalBuyPrice: totalBuyPrice.toString(),
        pilotCount: totalPilots,
        scannerCount,
        involvedCount: totalPilots, // All pilots are involved
        splitDetails: {
            "Scanner Payout": results.scannerPayout,
            "Pilot Payout": results.pilotPayout,
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
            toastr.success('Split saved.');
            window.location.href = "/loot-summary";
        } else {
            toastr.error('Error saving split.');
            console.error("Error saving loot split:", response.status, response.statusText);
        }
    } catch (error) {
        toastr.error('Error saving split.');
        console.error("Error saving loot split:", error);
    }
}


