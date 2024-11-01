// static/js/utils.js

export const noDataPlugin = {
    id: 'noData',
    afterDraw: function(chart) { // Changed from beforeDraw to afterDraw
        // Determine if the chart has no data
        const hasData = chart.data.datasets.some(dataset => {
            return dataset.data && dataset.data.length > 0 && dataset.data.some(value => {
                if (typeof value === 'object' && value !== null) {
                    return Object.values(value).some(val => val !== null && val !== undefined && val !== '');
                }
                return value !== null && value !== undefined && value !== '';
            });
        });

        if (!hasData) {
            // Retrieve the chart title
            const chartTitle = chart.options.plugins.title && chart.options.plugins.title.text ? chart.options.plugins.title.text : 'Unnamed Chart';

            // Log the chart title and data
            console.log(`No data for chart "${chartTitle}". Chart data:`, chart.data);

            const { ctx, width, height } = chart;
            ctx.save();
            ctx.textAlign = 'center';
            ctx.textBaseline = 'middle';
            ctx.font = '20px Montserrat, sans-serif';
            ctx.fillStyle = '#ff4d4d'; // Customize as needed

            // Calculate position below the title
            const titleHeight = chart.options.plugins.title && chart.options.plugins.title.display ? 40 : 0; // Approximate title height
            const messageY = height / 2 + titleHeight / 2;

            ctx.fillText('No data available', width / 2, messageY);
            ctx.restore();
        } else {
            const chartTitle = chart.options.plugins.title && chart.options.plugins.title.text ? chart.options.plugins.title.text : 'Unnamed Chart';

            console.log(`Data for chart "${chartTitle}". Chart data:`, chart.data);
        }
    }
};

/**
 * Generates common options for Chart.js charts.
 * @param {string} titleText - The title of the chart.
 * @param {Object} additionalOptions - Additional Chart.js options to merge.
 * @returns {Object} - The combined chart options.
 */
export function getCommonOptions(titleText, additionalOptions = {}) {
    // Destructure plugins and scales from additionalOptions to merge separately
    const { plugins: additionalPlugins = {}, scales: additionalScales = {}, datasets: additionalDatasets = {}, ...restOptions } = additionalOptions;

    return {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
            legend: { display: true, position: 'top', labels: { color: '#ffffff' } },
            title: {
                display: true,
                text: titleText,
                font: {
                    size: 18,
                    family: 'Montserrat, sans-serif',
                    weight: 'bold',
                },
                color: '#ffffff',
                align: 'center',
                padding: {
                    top: 10,
                    bottom: 30,
                },
                ...additionalPlugins.title,
            },
            tooltip: {
                mode: 'index',
                intersect: false,
                ...additionalPlugins.tooltip,
            },
            // Merge any other plugins here
            ...additionalPlugins,
        },
        scales: {
            x: {
                type: 'category',
                ticks: {
                    color: '#ffffff',
                    maxRotation: 45,
                    minRotation: 45,
                    autoSkip: false,
                    font: {
                        size: 10, // Reduced font size for better fit
                    },
                },
                grid: { display: false },
                title: {
                    display: true,
                    text: 'Characters',
                    color: '#ffffff',
                    font: {
                        size: 14,
                        family: 'Montserrat, sans-serif',
                        weight: 'bold',
                    },
                },
                stacked: false, // Set to true if you want stacked bars
                ...additionalScales.x,
            },
            y: {
                beginAtZero: true,
                ticks: { color: '#ffffff' },
                grid: { color: '#444' },
                title: {
                    display: true,
                    text: 'Count',
                    color: '#ffffff',
                    font: {
                        size: 14,
                        family: 'Montserrat, sans-serif',
                        weight: 'bold',
                    },
                },
                stacked: false, // Set to true if you want stacked bars
                ...additionalScales.y,
            },
            ...additionalScales,
        },
        layout: {
            padding: {
                left: 10,
                right: 10,
                top: 10,
                bottom: 10,
            },
            ...restOptions.layout,
        },
        interaction: {
            mode: 'index',
            intersect: false,
            ...additionalOptions.interaction,
        },
        // Merge dataset-specific options if needed
        datasets: {
            bar: {
                barPercentage: 0.5, // Adjust bar width (0 to 1)
                categoryPercentage: 0.5, // Adjust spacing between categories (0 to 1)
                ...additionalDatasets.bar,
            },
            ...additionalDatasets,
        },
        ...restOptions,
    };
}


/**
 * Truncates a label to a specified length and appends an ellipsis if necessary.
 * @param {string} label - The label to truncate.
 * @param {number} maxLength - The maximum length of the truncated label.
 * @returns {string} - The truncated label.
 */
export function truncateLabel(label, maxLength) {
    if (label.length > maxLength) {
        return label.substring(0, maxLength - 3) + '...';
    }
    return label;
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
