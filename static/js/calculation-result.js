function initCalculationResult() {
    const container = document.getElementById("calculation-result-container");
    container.innerHTML = `
        <div id="valuesContainer" class="bg-gray-800 p-4 rounded-lg text-center text-gray-100 space-y-2">
            <!-- Calculation results will be displayed here -->
        </div>
        <div id="splitRulesContainer" class="text-teal-200 text-lg mt-4">
            <!-- Split rules will be displayed here -->
        </div>`;
}

function calculateSplit(totalBuyPrice, totalPilots, scannerCount) {
    const results = {
        scannerPayout: 0,
        pilotPayout: 0,
        corpShare: 0,
        ruleExplanation: "",
    };

    if (isNaN(totalBuyPrice) || totalBuyPrice <= 0 || (totalPilots + scannerCount) <= 0) {
        results.ruleExplanation = "Invalid data provided.";
        return results;
    }

    const totalParticipants = totalPilots + scannerCount;
    const baseShare = Math.floor(totalBuyPrice / totalParticipants);

    // Calculate corporation share if base share exceeds threshold
    if (baseShare > 100_000_000) {
        results.corpShare = Math.floor(totalBuyPrice * 0.10); // Corporation gets 10%
        totalBuyPrice -= results.corpShare; // Adjust remaining loot
    }

    const basePayout = Math.floor(totalBuyPrice / totalParticipants);
    const scannerBonus = Math.floor(totalBuyPrice * 0.10); // Scanners get 10% bonus of total loot
    const scannerPayout = basePayout + Math.floor(scannerBonus / scannerCount); // Base + bonus per scanner
    const pilotPayout = basePayout; // Pilots get base payout
    const remainder = totalBuyPrice - ((scannerPayout * scannerCount) + (pilotPayout * totalPilots));

    // Add remainder to the first scanner payout
    results.scannerPayout = scannerPayout + (remainder > 0 ? remainder : 0);
    results.pilotPayout = pilotPayout;

    // Rule explanation
    results.ruleExplanation = results.corpShare > 0
        ? "10% Corp Share applied. Scanners receive 10% bonus on their base share."
        : "Scanners receive 10% bonus on their base share.";

    return results;
}


function calculateValues(totalBuyPrice) {
    const valuesContainer = document.getElementById("valuesContainer");
    const splitRulesContainer = document.getElementById("splitRulesContainer");

    const totalPilots = parseInt(document.getElementById("pilotCount").value) || 0;
    const scannerCount = parseInt(document.getElementById("scannerCount").value) || 0;

    valuesContainer.innerHTML = "";
    splitRulesContainer.innerHTML = "";

    // Call shared calculation function
    const results = calculateSplit(totalBuyPrice, totalPilots, scannerCount);

    if (results.ruleExplanation === "Invalid data provided.") {
        valuesContainer.innerHTML = "<p>Appraise Loot</p>";
        return;
    }

    // Display Scanner Payout
    if (scannerCount > 0) {
        valuesContainer.innerHTML += `
            <p>Scanner Payout: ${formatNumber(results.scannerPayout)} (each)
                <span class="clipboard-icon cursor-pointer text-gray-400 hover:text-green-500 ml-2" 
                      onclick="copyToClipboard(this, '${results.scannerPayout}')">
                    <i class="fas fa-clipboard"></i>
                </span>
            </p>`;
    }

    // Display Pilot Payout
    if (totalPilots > 0) {
        valuesContainer.innerHTML += `
            <p>Pilot Payout: ${formatNumber(results.pilotPayout)} (each)
                <span class="clipboard-icon cursor-pointer text-gray-400 hover:text-green-500 ml-2" 
                      onclick="copyToClipboard(this, '${results.pilotPayout}')">
                    <i class="fas fa-clipboard"></i>
                </span>
            </p>`;
    }

    // Display Corporation Share
    if (results.corpShare > 0) {
        valuesContainer.innerHTML += `
            <p>Corp Share: ${formatNumber(results.corpShare)}
                <span class="clipboard-icon cursor-pointer text-gray-400 hover:text-green-500 ml-2" 
                      onclick="copyToClipboard(this, '${results.corpShare}')">
                    <i class="fas fa-clipboard"></i>
                </span>
            </p>`;
    }

    // Display Rule Explanation
    splitRulesContainer.innerHTML = `<p class="text-center">${results.ruleExplanation}</p>`;
}



function recalculateSplit() {
    const totalBuyPrice = parseFloat(
        document.getElementById("lootInput")?.innerText.replace(/,/g, "") || "0"
    );

    const valuesContainer = document.getElementById("valuesContainer");
    const splitRulesContainer = document.getElementById("splitRulesContainer");

    if (!valuesContainer || !splitRulesContainer) {
        console.warn("Values or Split Rules container not found. Skipping split calculation.");
        return;
    }

    calculateValues(totalBuyPrice);
}

function copyToClipboard(element, value) {
    navigator.clipboard.writeText(value.toString()).then(() => {
        console.log("Copied to clipboard:", value);

        // Add temporary feedback styling
        element.classList.add('text-green-500', 'scale-125', 'transform', 'transition', 'duration-300');

        // Remove feedback after 1 second
        setTimeout(() => {
            element.classList.remove('text-green-500', 'scale-125');
        }, 1000);
    }).catch(err => {
        console.error("Failed to copy to clipboard", err);
    });
}
