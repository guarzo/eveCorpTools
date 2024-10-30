// static/js/utils.js

/**
 * Generates common Chart.js options with the ability to add a title.
 * @param {string} titleText - The text to display as the chart title.
 * @param {Object} [additionalOptions] - Additional Chart.js options to merge.
 * @returns {Object} - The Chart.js options object.
 */
export function getCommonOptions(titleText, additionalOptions = {}) {
    return {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
            legend: { display: true },
            tooltip: {
                callbacks: {
                    title: function (context) {
                        const index = context[0].dataIndex;
                        return context[0].chart.data.fullLabels
                            ? context[0].chart.data.fullLabels[index]
                            : context[0].chart.data.labels[index];
                    },
                },
            },
            title: {
                display: true,
                text: titleText,
                font: {
                    size: 18,
                    family: 'Montserrat, sans-serif',
                    weight: 'bold', // Changed from 'style' to 'weight'
                },
                color: '#ffffff',
                align: 'center',
                padding: {
                    top: 10,
                    bottom: 30,
                },
                // Allow overriding default title options
                ...(additionalOptions.plugins && additionalOptions.plugins.title ? { ...additionalOptions.plugins.title } : {}),
            },
            // Merge any additional plugin options provided
            ...(additionalOptions.plugins ? { ...additionalOptions.plugins } : {}),
        },
        scales: {
            x: {
                ticks: { color: '#ffffff' },
                grid: { color: '#444' },
                beginAtZero: true,
                // Allow overriding default x-axis options
                ...(additionalOptions.scales && additionalOptions.scales.x ? { ...additionalOptions.scales.x } : {}),
            },
            y: {
                ticks: {
                    color: '#ffffff',
                    autoSkip: false,
                },
                grid: { display: false },
                // Allow overriding default y-axis options
                ...(additionalOptions.scales && additionalOptions.scales.y ? { ...additionalOptions.scales.y } : {}),
            },
            // Merge any additional scales provided
            ...(additionalOptions.scales ? { ...additionalOptions.scales } : {}),
        },
        // Merge any additional options provided
        ...(additionalOptions ? { ...additionalOptions } : {}),
    };
}

/**
 * Truncates a label to a specified length and appends an ellipsis if necessary.
 * @param {string} label - The label to truncate.
 * @param {number} length - The maximum length of the truncated label.
 * @returns {string} - The truncated label.
 */
export function truncateLabel(label, length) {
    if (!label || typeof label !== 'string') {
        return '';
    }
    return label.length > length ? label.substring(0, length) + '...' : label;
}

/**
 * Returns a color from a predefined palette based on the index.
 * @param {number} index - The index to determine the color.
 * @returns {string} - The corresponding color in HEX format.
 */
export function getColor(index) {
    const colors = [
        '#FF6384', '#36A2EB', '#FFCE56', '#4BC0C0',
        '#9966FF', '#FF9F40', '#E7E9ED', '#76D7C4',
        '#C0392B', '#8E44AD', '#2ECC71', '#1ABC9C',
        '#3498DB', '#F1C40F', '#E67E22', '#95A5A6',
    ];
    return colors[index % colors.length];
}

/**
 * Validates the chart data to ensure it meets the required format.
 * @param {Array} data - The data to validate.
 * @param {string} chartName - The name of the chart (for logging purposes).
 * @returns {boolean} - Returns true if data is valid, false otherwise.
 */
export function validateChartData(data, chartName) {
    if (!Array.isArray(data)) {
        console.warn(`${chartName}: Data should be an array.`);
        return false;
    }

    if (data.length === 0) {
        console.warn(`${chartName}: Data array is empty.`);
        return false;
    }

    // Additional validation logic can be added here based on chart requirements
    // For example, checking for required fields in each data object

    return true;
}
