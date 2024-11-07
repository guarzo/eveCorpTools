// loot-summary.js

document.addEventListener("DOMContentLoaded", function() {
    fetchLootSummaries();
});

async function fetchLootSummaries() {
    try {
        const response = await fetch('/fetch-loot-splits');

        // Check if the response is OK (status code 200-299)
        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`);
        }

        // Attempt to parse the response as JSON
        const data = await response.json();

        console.log('Fetched data:', data); // Debugging line

        // Check if data is an array
        if (Array.isArray(data)) {
            if (data.length > 0) {
                // Process data
                const processedData = data.map((row, index) => ({
                    ...row,
                    id: index,
                    date: row.date
                        ? luxon.DateTime.fromISO(row.date, { zone: 'utc' }).toFormat('yyyy-MM-dd HH:mm:ss')
                        : 'N/A',
                    totalBuyPrice: row.totalBuyPrice
                        ? formatNumber(Number(row.totalBuyPrice))
                        : '0', // Convert to number
                }));


                window.lootSummaryTable = new Tabulator("#lootSummaryTable", {
                    data: processedData,
                    layout: "fitData",
                    responsiveLayout: "hide",
                    columns: [
                        { title: "Date", field: "date", minWidth: 300 },
                        { title: "Battle Report", field: "battleReport", formatter: "link", formatterParams: { labelField: "battleReport", urlField: "battleReport", minWidth: 300 } },
                        { title: "Total Buy Price", field: "totalBuyPrice", minWidth: 300 },
                        {
                            title: "Remove",
                            formatter: "buttonCross",
                            hozAlign: "center",
                            headerSort: false,
                            minWidth: 50,
                            cellClick: function (e, cell) {
                                const rowData = cell.getRow().getData();
                                confirmDelete(rowData.id); // Pass row id to confirmDelete
                            }
                        }
                    ],
                    // Additional Tabulator configurations as needed
                });

                // Set up row click event
                window.lootSummaryTable.on("rowClick", function(e, row) {
                    const details = row.getData();
                    document.getElementById("selectedRowId").value = details.id;
                    document.getElementById("detailDate").innerText = details.date || 'N/A';
                    document.getElementById("detailBattleReport").innerText = details.battleReport || 'N/A';
                    document.getElementById("detailBattleReport").href = details.battleReport || '#';
                    document.getElementById("detailTotalBuyPrice").innerText = details.totalBuyPrice || '0';
                    displaySplitDetails(details.splitDetails);
                    document.getElementById("detailContainer").style.display = "block";
                });

                // Ensure the lootSummaryTable is visible
                document.getElementById("lootSummaryTable").style.display = "block";

                // Hide the 'no data' message if it exists
                const noDataMessage = document.getElementById("noDataMessage");
                if (noDataMessage) {
                    noDataMessage.style.display = "none";
                }

            } else {
                // No data available
                displayNoDataMessage();
            }
        } else {
            // Unexpected data format
            console.warn("Unexpected data format:", data);
            displayErrorMessage("Received data in unexpected format.");
        }

    } catch (error) {
        console.error("Error fetching loot summaries:", error);
        displayErrorMessage(error.message || "An unexpected error occurred.");
    }
}

function confirmDelete(id) {
    Swal.fire({
        title: 'Are you sure?',
        text: "This action cannot be undone!",
        icon: 'warning',
        showCancelButton: true,
        confirmButtonColor: '#d33',
        cancelButtonColor: '#3085d6',
        confirmButtonText: 'Yes, delete it!',
        cancelButtonText: 'Cancel',
        reverseButtons: false
    }).then((result) => {
        if (result.isConfirmed) {
            deleteDetails(id); // Pass the id to deleteDetails
            Swal.fire({
                title: 'Deleted!',
                text: 'The loot split entry has been deleted.',
                icon: 'success',
                confirmButtonColor: '#3085d6'
            });
        }
    });
}

function deleteDetails(id) {
    fetch('/delete-loot-split', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ id: id })
    })
        .then(response => {
            if (response.ok) {
                console.log("Loot split deleted successfully.");
                fetchLootSummaries(); // Refresh the table
                document.getElementById("detailContainer").style.display = "none"; // Hide details container
            } else {
                console.error("Error deleting loot split.");
                toastr.error('Failed to delete loot split.');
            }
        })
        .catch(error => {
            console.error("Error deleting loot split.", error);
            toastr.error('An error occurred while deleting the loot split.');
        });
}

function formatNumber(num) {
    if (typeof num !== 'number') return '0';
    return num.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",");
}

function displaySplitDetails(splitDetails) {
    const splitDetailsContainer = document.getElementById("splitDetails");
    splitDetailsContainer.innerHTML = '';
    if (splitDetails && typeof splitDetails === 'object' && Object.keys(splitDetails).length > 0) {
        for (const key in splitDetails) {
            if (splitDetails.hasOwnProperty(key)) {
                const value = splitDetails[key];
                splitDetailsContainer.innerHTML += `
                    <p><strong>${key}:</strong> ${formatNumber(value)}
                        <span class="clipboard-icon cursor-pointer text-gray-400 hover:text-green-500 ml-2" onclick="copyToClipboard(this, '${value}')">
                            <i class="fas fa-clipboard"></i>
                        </span>
                    </p>`;
            }
        }
    } else {
        splitDetailsContainer.innerHTML = '<p>No split details available</p>';
    }
}

// Function to copy text to clipboard and add temporary visual feedback
function copyToClipboard(element, value) {
    navigator.clipboard.writeText(value.toString()).then(() => {
        console.log("Copied to clipboard:", value);

        // Visual feedback using Tailwind classes
        element.classList.add('text-green-500', 'scale-125', 'transform', 'transition', 'duration-300');

        // Remove feedback after 1 second
        setTimeout(() => {
            element.classList.remove('text-green-500', 'scale-125');
        }, 1000);
    }).catch(err => {
        console.error("Failed to copy to clipboard", err);
    });
}

function confirmDelete() {
    Swal.fire({
        title: 'Are you sure?',
        text: "This action cannot be undone!",
        icon: 'warning',
        showCancelButton: true,
        confirmButtonColor: '#d33',
        cancelButtonColor: '#3085d6',
        confirmButtonText: 'Yes, delete it!',
        cancelButtonText: 'Cancel',
        reverseButtons: false
    }).then((result) => {
        if (result.isConfirmed) {
            deleteDetails();
            Swal.fire({
                title: 'Deleted!',
                text: 'The loot split entry has been deleted.',
                icon: 'success',
                confirmButtonColor: '#3085d6'
            });
        }
    });
}

function deleteDetails() {
    const id = parseInt(document.getElementById("selectedRowId").value);
    fetch('/delete-loot-split', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ id: id })
    })
        .then(response => {
            if (response.ok) {
                console.log("Loot split deleted successfully.");
                fetchLootSummaries(); // Refresh the table
                document.getElementById("detailContainer").style.display = "none"; // Hide details container
            } else {
                console.error("Error deleting loot split.");
                toastr.error('Failed to delete loot split.');
            }
        })
        .catch(error => {
            console.error("Error deleting loot split.", error);
            toastr.error('An error occurred while deleting the loot split.');
        });
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
