function initCalculationResult() {
    const container = document.getElementById("calculation-result-container");
    container.innerHTML = `
        <div id="valuesContainer" class="highlight">
            <!-- Calculation results will be displayed here -->
        </div>
        <div id="splitRulesContainer" class="highlight center-text">
            <!-- Split rules will be displayed here -->
        </div>`;
}

function calculateValues(totalBuyPrice) {
    const valuesContainer = document.getElementById('valuesContainer');
    const splitRulesContainer = document.getElementById('splitRulesContainer');
    if (!valuesContainer || !splitRulesContainer) {
        console.error('valuesContainer or splitRulesContainer not found');
        return;
    }

    const scoutBox = document.getElementById('scout');
    const involvedBox = document.getElementById('involved');
    const scout = scoutBox.querySelector('.draggable-box');
    const involved = involvedBox.querySelectorAll('.draggable-box');

    valuesContainer.innerHTML = '';
    splitRulesContainer.innerHTML = '';

    if (isNaN(totalBuyPrice) || totalBuyPrice <= 0) {
        valuesContainer.innerHTML = ''; // Clear the container
        splitRulesContainer.innerHTML = ''; // Clear the rules container
        return;
    }

    let splitRuleText = '';
    let shares = [];
    let corpShare = 0;

    if (scout && involved.length + 1 < 3) {
        const equalShare = Math.floor(totalBuyPrice / (involved.length + 1));
        const remainder = totalBuyPrice - (equalShare * (involved.length + 1));
        shares.push(equalShare + remainder);

        valuesContainer.innerHTML += `<p>Scout (${scout.innerHTML}): ${formatNumber(equalShare + remainder)}<span class="clipboard-icon" onclick="copyToClipboard(this, '${equalShare + remainder}')">📋</span></p>`;
        involved.forEach(box => {
            valuesContainer.innerHTML += `<p>Involved (${box.innerHTML}): ${formatNumber(equalShare)}<span class="clipboard-icon" onclick="copyToClipboard(this, '${equalShare}')">📋</span></p>`;
            shares.push(equalShare);
        });

        splitRuleText = 'Even Split';
    } else if (scout && involved.length === 0) {
        const scoutValue = Math.floor(totalBuyPrice * 0.90);
        corpShare = Math.floor(totalBuyPrice * 0.10);
        const remainder = totalBuyPrice - (scoutValue + corpShare);
        corpShare += remainder;
        shares.push(scoutValue);

        valuesContainer.innerHTML += `<p>Scout (${scout.innerHTML}): ${formatNumber(scoutValue)}<span class="clipboard-icon" onclick="copyToClipboard(this, '${scoutValue}')">📋</span></p>`;
        valuesContainer.innerHTML += `<p>Corp: ${formatNumber(corpShare)}<span class="clipboard-icon" onclick="copyToClipboard(this, '${corpShare}')">📋</span></p>`;

        splitRuleText = '90% to Scout, 10% to Corp';
    } else if (scout && involved.length > 0) {
        const scoutValue = Math.floor(totalBuyPrice * 0.10);
        corpShare = Math.floor(totalBuyPrice * 0.10);
        const involvedValue = Math.floor((totalBuyPrice * 0.80) / (involved.length + 1)); // including scout
        const remainder = totalBuyPrice - (scoutValue + corpShare + (involvedValue * (involved.length + 1)));
        corpShare += remainder;

        const totalScoutValue = scoutValue + involvedValue;
        shares.push(totalScoutValue);

        valuesContainer.innerHTML += `<p>Scout (${scout.innerHTML}): ${formatNumber(totalScoutValue)}<span class="clipboard-icon" onclick="copyToClipboard(this, '${totalScoutValue}')">📋</span></p>`;
        involved.forEach(box => {
            valuesContainer.innerHTML += `<p>Involved (${box.innerHTML}): ${formatNumber(involvedValue)}<span class="clipboard-icon" onclick="copyToClipboard(this, '${involvedValue}')">📋</span></p>`;
            shares.push(involvedValue);
        });

        valuesContainer.innerHTML += `<p>Corp: ${formatNumber(corpShare)}<span class="clipboard-icon" onclick="copyToClipboard(this, '${corpShare}')">📋</span></p>`;

        splitRuleText = '10% to the Scout, 80% to involved + scout, 10% to corp';
    } else if (!scout && involved.length < 3) {
        const equalShare = Math.floor(totalBuyPrice / involved.length);
        const remainder = totalBuyPrice - (equalShare * involved.length);

        involved.forEach((box, index) => {
            const adjustedInvolvedValue = index === involved.length - 1 ? equalShare + remainder : equalShare;
            valuesContainer.innerHTML += `<p>Involved (${box.innerHTML}): ${formatNumber(adjustedInvolvedValue)}<span class="clipboard-icon" onclick="copyToClipboard(this, '${adjustedInvolvedValue}')">📋</span></p>`;
            shares.push(adjustedInvolvedValue);
        });

        splitRuleText = 'Even Split';
    } else if (!scout && involved.length >= 3) {
        corpShare = Math.floor(totalBuyPrice * 0.10);
        const involvedValue = Math.floor(totalBuyPrice * 0.90 / involved.length);
        const remainder = totalBuyPrice - (corpShare + involvedValue * involved.length);

        valuesContainer.innerHTML += `<p>Corp: ${formatNumber(corpShare)}<span class="clipboard-icon" onclick="copyToClipboard(this, '${corpShare}')">📋</span></p>`;
        involved.forEach((box, index) => {
            const adjustedInvolvedValue = index === involved.length - 1 ? involvedValue + remainder : involvedValue;
            valuesContainer.innerHTML += `<p>Involved (${box.innerHTML}): ${formatNumber(adjustedInvolvedValue)}<span class="clipboard-icon" onclick="copyToClipboard(this, '${adjustedInvolvedValue}')">📋</span></p>`;
            shares.push(adjustedInvolvedValue);
        });

        splitRuleText = '90% split to involved, 10% to Corp';
    } else {
        corpShare = Math.floor(totalBuyPrice * 0.10);
        const remainder = totalBuyPrice - corpShare;
        corpShare += remainder;
        shares.push(corpShare);

        valuesContainer.innerHTML += `<p>Corp: ${formatNumber(corpShare)}<span class="clipboard-icon" onclick="copyToClipboard(this, '${corpShare}')">📋</span></p>`;
        splitRuleText = '100% to the Corp';
    }

    // Check if any pilot receives less than 100,000,000 ISK (excluding corp share)
    const pilotShares = shares.filter(share => share !== corpShare);
    if (Math.min(...pilotShares) < 100000000) {
        console.log('Total Buy Price:', totalBuyPrice);
        console.log('Original split before adjustment:', shares);
        valuesContainer.innerHTML = ''; // Clear the current split details

        if (scout) {
            const scoutValue = Math.floor(totalBuyPrice * 0.25);
            corpShare = Math.floor(totalBuyPrice * 0.75);
            const remainder = totalBuyPrice - (scoutValue + corpShare);
            corpShare += remainder;

            valuesContainer.innerHTML += `<p>Scout (${scout.innerHTML}): ${formatNumber(scoutValue)}<span class="clipboard-icon" onclick="copyToClipboard(this, '${scoutValue}')">📋</span></p>`;
            valuesContainer.innerHTML += `<p>Corp: ${formatNumber(corpShare)}<span class="clipboard-icon" onclick="copyToClipboard(this, '${corpShare}')">📋</span></p>`;
            splitRuleText = '25% to the Scout, 75% to the Corp';
        } else {
            corpShare = Math.floor(totalBuyPrice);
            valuesContainer.innerHTML += `<p>Corp: ${formatNumber(corpShare)}<span class="clipboard-icon" onclick="copyToClipboard(this, '${corpShare}')">📋</span></p>`;
            splitRuleText = '100% to the Corp';
        }
    }

    splitRulesContainer.innerHTML = `<p>${splitRuleText}</p>`;

    valuesContainer.classList.add('highlight');
    setTimeout(() => valuesContainer.classList.remove('highlight'), 2000);
}

// Function to copy text to clipboard and add a temporary visual indicator
function copyToClipboard(icon, value) {
    const tempInput = document.createElement("input");
    tempInput.value = value;
    document.body.appendChild(tempInput);
    tempInput.select();
    document.execCommand("copy");
    document.body.removeChild(tempInput);

    // Temporarily change icon color and size to indicate copy success
    icon.classList.add("copied");
    requestAnimationFrame(() => {
        setTimeout(() => {
            icon.classList.remove("copied");
        }, 1000);
    });
}
