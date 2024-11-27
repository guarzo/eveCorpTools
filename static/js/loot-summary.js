// loot-summary.js

document.addEventListener("DOMContentLoaded", function() {
    fetchLootSummaries();
});

async function fetchLootSummaries() {
    try {
        const response = await fetch('/fetch-loot-splits');
        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`);
        }

        const data = await response.json();
        console.log('Fetched data:', data);

        if (Array.isArray(data) && data.length > 0) {
            // Map data and preserve backend IDs
            const processedData = data.map((row) => ({
                ...row,
                date: row.date
                    ? luxon.DateTime.fromISO(row.date, { zone: 'utc' }).toFormat('yyyy-MM-dd HH:mm:ss')
                    : 'N/A',
                totalBuyPrice: row.totalBuyPrice
                    ? formatNumber(Number(row.totalBuyPrice))
                    : '0',
            }));

            window.lootSummaryTable = new Tabulator("#lootSummaryTable", {
                theme: "midnight",
                data: processedData,
                layout: "fitColumns",
                responsiveLayout: "hide",
                columns: [
                    { title: `<i class="fas fa-calendar-alt text-teal-500"></i> Date`, field: "date", minWidth: 150, sorter: "datetime" },
                    {
                        title: `<i class="fas fa-link text-teal-500"></i> Battle Report`,
                        field: "battleReport",
                        headerTooltip: "Link to the battle report or description",
                        formatter: function (cell) {
                            const value = cell.getValue();
                            if (value.startsWith('http://') || value.startsWith('https://')) {
                                return `<a href="${value}" target="_blank" class="text-teal-300 hover:text-teal-500" title="View Battle Report">
                                    <i class="fas fa-external-link-alt"></i>
                                </a>`;
                            }
                            return `<span class="text-gray-200">${value}</span>`;
                        },
                        minWidth: 50,
                    },
                    {
                        title: `<i class="fas fa-coins text-yellow-500"></i> Total Buy Price`,
                        field: "totalBuyPrice",
                        minWidth: 150,
                        formatter: (cell) => `<i class="fas fa-coins text-yellow-500 mr-2"></i>${cell.getValue()} ISK`,
                        headerTooltip: "Total value of the loot",
                        sorter: "number",
                    },
                ],
                rowFormatter: (row) => row.getElement().style.cursor = "pointer",
            });

            window.lootSummaryTable.on("rowClick", function (e, row) {
                const details = row.getData();

                // Use the correct field name for ID
                document.getElementById("selectedRowId").value = details.id; // Backend-provided id
                document.getElementById("detailDate").innerText = details.date || 'N/A';
                document.getElementById("detailBattleReportInput").value = details.battleReport || '';
                document.getElementById("detailTotalBuyPrice").innerText = details.totalBuyPrice || '0';

                const battleReportLabel = document.getElementById("detailBattleReportLabel");
                if (details.battleReport.startsWith('http://') || details.battleReport.startsWith('https://')) {
                    battleReportLabel.href = details.battleReport;
                    battleReportLabel.target = "_blank";
                    battleReportLabel.className = "text-teal-300 hover:underline";
                } else {
                    battleReportLabel.href = "#";
                    battleReportLabel.className = "text-gray-400 cursor-default";
                }

                displaySplitDetails(details.splitDetails);
                openDetailModal();
            });


            document.getElementById("lootSummaryTable").style.display = "block";
            const noDataMessage = document.getElementById("noDataMessage");
            if (noDataMessage) {
                noDataMessage.style.display = "none";
            }
        } else {
            displayNoDataMessage();
        }
    } catch (error) {
        console.error("Error fetching loot summaries:", error);
        displayErrorMessage(error.message || "An unexpected error occurred.");
    }
}

function saveBattleReportUpdate() {
    const id = parseInt(document.getElementById("selectedRowId").value, 10);
    const battleReport = document.getElementById("detailBattleReportInput").value;

    if (isNaN(id) || id <= 0) {
        toastr.error('No loot split selected for updating.');
        console.error('Invalid ID for update:', id); // Debugging
        return;
    }

    fetch('/update-loot-split', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ id, battleReport }),
    })
        .then((response) => {
            if (response.ok) {
                toastr.success('Battle report updated successfully.');
                fetchLootSummaries();
                closeDetailModal();
            } else {
                toastr.error('Failed to update battle report.');
                console.error(`Failed to update battle report for ID: ${id}`);
            }
        })
        .catch((error) => {
            console.error('Error updating battle report:', error);
            toastr.error('An error occurred while updating the battle report.');
        });
}



function confirmDelete() {
    const id = parseInt(document.getElementById("selectedRowId").value, 10);

    if (isNaN(id) || id <= 0) {
        toastr.error('No loot split selected for deletion.');
        console.error('Invalid ID for deletion:', id); // Debugging
        return;
    }

    Swal.fire({
        title: 'Are you sure?',
        text: "This action cannot be undone!",
        icon: 'warning',
        showCancelButton: true,
        confirmButtonColor: '#d33',
        cancelButtonColor: '#3085d6',
        confirmButtonText: 'Yes, delete it!',
        cancelButtonText: 'Cancel',
    }).then((result) => {
        if (result.isConfirmed) {
            deleteDetails(id);
        }
    });
}


function deleteDetails(id) {
    fetch('/delete-loot-split', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ id }),
    })
        .then((response) => {
            if (response.ok) {
                toastr.success('Loot split deleted successfully.');
                fetchLootSummaries(); // Refresh the table
                closeDetailModal(); // Close the modal
            } else {
                toastr.error('Failed to delete loot split.');
                console.error(`Failed to delete loot split with ID: ${id}`);
            }
        })
        .catch((error) => {
            console.error('Error deleting loot split:', error);
            toastr.error('An error occurred while deleting the loot split.');
        });
}

// Function to open the detail modal
function openDetailModal() {
    document.getElementById("detailModal").classList.remove('hidden');
}

// Function to close the detail modal
function closeDetailModal() {
    document.getElementById("detailModal").classList.add('hidden');
}

function formatNumber(num) {
    if (typeof num !== 'number') return '0';
    return num.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",");
}

function displaySplitDetails(splitDetails) {
    const splitDetailsContainer = document.getElementById("splitDetails");
    splitDetailsContainer.innerHTML = '';
    if (splitDetails && typeof splitDetails === 'object' && Object.keys(splitDetails).length > 0) {
        const table = document.createElement('table');
        table.classList.add('w-full', 'text-left', 'mt-4');
        for (const key in splitDetails) {
            if (splitDetails.hasOwnProperty(key)) {
                const value = splitDetails[key];
                const row = document.createElement('tr');

                const cellKey = document.createElement('td');
                cellKey.classList.add('py-2', 'font-semibold', 'text-gray-200');
                cellKey.innerText = key;

                const cellValue = document.createElement('td');
                cellValue.classList.add('py-2', 'text-teal-200', 'text-right');
                cellValue.innerHTML = `${formatNumber(Number(value))} ISK
                    <span class="clipboard-icon cursor-pointer text-gray-400 hover:text-green-500 ml-2" onclick="copyToClipboard(this, '${value}')">
                        <i class="fas fa-clipboard"></i>
                    </span>`;

                row.appendChild(cellKey);
                row.appendChild(cellValue);
                table.appendChild(row);
            }
        }
        splitDetailsContainer.appendChild(table);
    } else {
        splitDetailsContainer.innerHTML = '<p>No split details available</p>';
    }
}

// Function to display a message when no data is available
function displayNoDataMessage() {
    // Hide the table
    document.getElementById("lootSummaryTable").style.display = "none";

    // Create and display the message
    let message = document.getElementById("noDataMessage");
    if (!message) {
        message = document.createElement("div");
        message.id = "noDataMessage";
        message.className = "text-center text-xl text-yellow-400 mt-8";
        message.innerText = "No saved loot splits available.";
        document.querySelector("main").appendChild(message);
    } else {
        message.style.display = "block";
    }

    // Hide the detail container if visible
    document.getElementById("detailContainer").style.display = "none";
}

// Function to display an error message
function displayErrorMessage(errorMessage) {
    // Hide the table and no data message
    document.getElementById("lootSummaryTable").style.display = "none";
    const noDataMessage = document.getElementById("noDataMessage");
    if (noDataMessage) {
        noDataMessage.style.display = "none";
    }

    // Create and display the error message
    let errorMsg = document.getElementById("errorMessage");
    if (!errorMsg) {
        errorMsg = document.createElement("div");
        errorMsg.id = "errorMessage";
        errorMsg.className = "text-center text-xl text-red-500 mt-8";
        errorMsg.innerText = "Failed to load loot summaries. Please try again later.";
        document.querySelector("main").appendChild(errorMsg);
    } else {
        errorMsg.style.display = "block";
    }

    // Hide the detail container if visible
    document.getElementById("detailContainer").style.display = "none";
}
