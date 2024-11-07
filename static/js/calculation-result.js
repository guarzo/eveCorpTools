function initCalculationResult() {
    const container = document.getElementById("calculation-result-container");
    container.innerHTML = `
        <div id="valuesContainer" class="bg-gray-800 p-4 rounded-lg text-center space-y-2">
            <!-- Calculation results will be displayed here -->
        </div>
        <div id="splitRulesContainer" class="text-yellow-400 text-lg mt-4">
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
    let corpShare = 0;

    if (scout && involved.length + 1 < 3) {
        const equalShare = Math.floor(totalBuyPrice / (involved.length + 1));
        const remainder = totalBuyPrice - (equalShare * (involved.length + 1));
        valuesContainer.innerHTML += `<p>Scout: ${scout.textContent} - ${formatNumber(equalShare + remainder)}
            <span class="clipboard-icon cursor-pointer text-gray-400 hover:text-green-500 ml-2" onclick="copyToClipboard(this, '${equalShare + remainder}')">
                <i class="fas fa-clipboard"></i>
            </span>
        </p>`;

        involved.forEach(box => {
            valuesContainer.innerHTML += `<p>Involved: ${box.textContent} - ${formatNumber(equalShare)}
                <span class="clipboard-icon cursor-pointer text-gray-400 hover:text-green-500 ml-2" onclick="copyToClipboard(this, '${equalShare}')">
                    <i class="fas fa-clipboard"></i>
                </span>
            </p>`;
        });
        splitRuleText = 'Even Split';
    } else if (scout && involved.length === 0) {
        const scoutValue = Math.floor(totalBuyPrice * 0.90);
        corpShare = Math.floor(totalBuyPrice * 0.10);
        const remainder = totalBuyPrice - (scoutValue + corpShare);
        corpShare += remainder;

        valuesContainer.innerHTML += `<p>Scout: ${scout.textContent} - ${formatNumber(scoutValue)}
            <span class="clipboard-icon cursor-pointer text-gray-400 hover:text-green-500 ml-2" onclick="copyToClipboard(this, '${scoutValue}')">
                <i class="fas fa-clipboard"></i>
            </span>
        </p>`;
        valuesContainer.innerHTML += `<p>Corp: ${formatNumber(corpShare)}
            <span class="clipboard-icon cursor-pointer text-gray-400 hover:text-green-500 ml-2" onclick="copyToClipboard(this, '${corpShare}')">
                <i class="fas fa-clipboard"></i>
            </span>
        </p>`;
        splitRuleText = '90% to Scout, 10% to Corp';
    } else if (scout && involved.length > 0) {
        const scoutValue = Math.floor(totalBuyPrice * 0.10);
        corpShare = Math.floor(totalBuyPrice * 0.10);
        const involvedValue = Math.floor((totalBuyPrice * 0.80) / (involved.length + 1)); // including scout
        const remainder = totalBuyPrice - (scoutValue + corpShare + (involvedValue * (involved.length + 1)));
        corpShare += remainder;

        const totalScoutValue = scoutValue + involvedValue;

        valuesContainer.innerHTML += `<p>Scout: ${scout.textContent} - ${formatNumber(totalScoutValue)}
            <span class="clipboard-icon cursor-pointer text-gray-400 hover:text-green-500 ml-2" onclick="copyToClipboard(this, '${totalScoutValue}')">
                <i class="fas fa-clipboard"></i>
            </span>
        </p>`;
        involved.forEach(box => {
            valuesContainer.innerHTML += `<p>Involved: ${box.textContent} - ${formatNumber(involvedValue)}
                <span class="clipboard-icon cursor-pointer text-gray-400 hover:text-green-500 ml-2" onclick="copyToClipboard(this, '${involvedValue}')">
                    <i class="fas fa-clipboard"></i>
                </span>
            </p>`;
        });

        valuesContainer.innerHTML += `<p>Corp: ${formatNumber(corpShare)}
            <span class="clipboard-icon cursor-pointer text-gray-400 hover:text-green-500 ml-2" onclick="copyToClipboard(this, '${corpShare}')">
                <i class="fas fa-clipboard"></i>
            </span>
        </p>`;
        splitRuleText = '10% to Scout, 80% to involved + scout, 10% to Corp';
    } else if (!scout && involved.length < 3) {
        const equalShare = Math.floor(totalBuyPrice / involved.length);
        const remainder = totalBuyPrice - (equalShare * involved.length);

        involved.forEach((box, index) => {
            const adjustedInvolvedValue = index === involved.length - 1 ? equalShare + remainder : equalShare;
            valuesContainer.innerHTML += `<p>Involved: ${box.textContent} - ${formatNumber(adjustedInvolvedValue)}
                <span class="clipboard-icon cursor-pointer text-gray-400 hover:text-green-500 ml-2" onclick="copyToClipboard(this, '${adjustedInvolvedValue}')">
                    <i class="fas fa-clipboard"></i>
                </span>
            </p>`;
        });
        splitRuleText = 'Even Split';
    } else if (!scout && involved.length >= 3) {
        corpShare = Math.floor(totalBuyPrice * 0.10);
        const involvedValue = Math.floor(totalBuyPrice * 0.90 / involved.length);
        const remainder = totalBuyPrice - (corpShare + involvedValue * involved.length);

        valuesContainer.innerHTML += `<p>Corp: ${formatNumber(corpShare)}
            <span class="clipboard-icon cursor-pointer text-gray-400 hover:text-green-500 ml-2" onclick="copyToClipboard(this, '${corpShare}')">
                <i class="fas fa-clipboard"></i>
            </span>
        </p>`;
        involved.forEach((box, index) => {
            const adjustedInvolvedValue = index === involved.length - 1 ? involvedValue + remainder : involvedValue;
            valuesContainer.innerHTML += `<p>Involved: ${box.textContent} - ${formatNumber(adjustedInvolvedValue)}
                <span class="clipboard-icon cursor-pointer text-gray-400 hover:text-green-500 ml-2" onclick="copyToClipboard(this, '${adjustedInvolvedValue}')">
                    <i class="fas fa-clipboard"></i>
                </span>
            </p>`;
        });
        splitRuleText = '90% split to involved, 10% to Corp';
    } else {
        corpShare = Math.floor(totalBuyPrice * 0.10);
        const remainder = totalBuyPrice - corpShare;
        corpShare += remainder;

        valuesContainer.innerHTML += `<p>Corp: ${formatNumber(corpShare)}
            <span class="clipboard-icon cursor-pointer text-gray-400 hover:text-green-500 ml-2" onclick="copyToClipboard(this, '${corpShare}')">
                <i class="fas fa-clipboard"></i>
            </span>
        </p>`;
        splitRuleText = '100% to the Corp';
    }

    splitRulesContainer.innerHTML = `<p>Split Rule Used: ${splitRuleText}</p>`;

    valuesContainer.classList.add('highlight');
    setTimeout(() => valuesContainer.classList.remove('highlight'), 2000);
}

function copyToClipboard(element, value) {
    navigator.clipboard.writeText(value).then(() => {
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
