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
 * Truncates a label to a specified maximum length, adding ellipsis if necessary.
 * @param {string} label - The original label.
 * @param {number} maxLength - The maximum allowed length.
 * @returns {string} - The truncated label.
 */
export function truncateLabel(label, maxLength) {
    if (label.length > maxLength) {
        return label.substring(0, maxLength - 3) + '...';
    }
    return label;
}

/**
 * Validates chart data to ensure it meets required criteria.
 * @param {Array} data - The data array to validate.
 * @param {string} chartName - The name of the chart for logging purposes.
 * @returns {boolean} - Returns true if data is valid, false otherwise.
 */
export function validateChartDataArray(data, chartName) {
    if (!Array.isArray(data) || data.length === 0) {
        console.warn(`No data available for chart "${chartName}".`);
        return false;
    }
    return true;
}

export function validateOurShipsUsedData(data, chartName) {
    if (typeof data !== 'object' || data === null) {
        console.warn(`Invalid data format for chart "${chartName}". Expected an object.`);
        return false;
    }

    const requiredKeys = ['Characters', 'ShipNames', 'SeriesData'];
    const hasAllKeys = requiredKeys.every(key => key in data);
    if (!hasAllKeys) {
        console.warn(`Incomplete data for chart "${chartName}". Missing keys: ${requiredKeys.filter(key => !(key in data)).join(', ')}`);
        return false;
    }

    if (!Array.isArray(data.Characters) || data.Characters.length === 0) {
        console.warn(`No characters data available for chart "${chartName}".`);
        return false;
    }
    if (!Array.isArray(data.ShipNames) || data.ShipNames.length === 0) {
        console.warn(`No ship names data available for chart "${chartName}".`);
        return false;
    }
    if (typeof data.SeriesData !== 'object' || Object.keys(data.SeriesData).length === 0) {
        console.warn(`No series data available for chart "${chartName}".`);
        return false;
    }

    return true;
}

/**
 * Generates common options for Chart.js charts.
 * @param {string} titleText - The title of the chart.
 * @param {Object} additionalOptions - Additional Chart.js options to merge.
 * @returns {Object} - The combined chart options.
 */

// static/js/utils.js

export function getCommonOptions(titleText, additionalOptions = {}) {
    const {
        plugins: additionalPlugins = {},
        scales: additionalScales = {},
        datasets: additionalDatasets = {},
        ...restOptions
    } = additionalOptions;

    return {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
            noData: {
                text: 'No data available for this chart',
                color: '#ffffff',
                font: {
                    size: 20,
                    family: 'Montserrat, sans-serif',
                    weight: 'bold',
                },
            },
            legend: {
                display: true,
                position: 'top',
                labels: { color: '#ffffff', font: { size: 12 } }
            },
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
                ...additionalPlugins.title, // Merge any additional title options
            },
            tooltip: {
                mode: 'index',
                intersect: false,
                callbacks: {
                    label: function (context) {
                        const dataset = context.dataset;
                        const shipName = dataset.label || '';
                        const value = context.parsed.x !== undefined ? context.parsed.x : context.parsed.y;
                        const index = context.dataIndex;
                        const total = dataset.percentage ? dataset.percentage[index] : 1;
                        const percentage = ((value / total) * 100).toFixed(2);
                        return `${shipName}: ${value} (${percentage}%)`;
                    },
                },
                ...additionalPlugins.tooltip, // Merge any additional tooltip options
            },
            datalabels: {
                color: '#ffffff',
                anchor: 'end',
                align: 'right',
                formatter: (value, context) => {
                    const dataset = context.dataset;
                    const index = context.dataIndex;
                    const total = dataset.percentage ? dataset.percentage[index] : 1;
                    const percentage = ((value / total) * 100).toFixed(1);
                    return `${value} (${percentage}%)`;
                },
                font: {
                    size: 10,
                    weight: 'bold',
                },
                ...additionalPlugins.datalabels, // Merge any additional datalabels options
            },
            // Merge any other plugins here without duplicating the 'plugins' property
            ...additionalPlugins,
        },
        scales: {
            x: {
                stacked: true,
                ticks: { color: '#ffffff' },
                grid: { display: false },
                ...additionalScales.x, // Merge any additional x-axis options
            },
            y: {
                stacked: true,
                ticks: {
                    color: '#ffffff',
                    autoSkip: false,
                },
                grid: { display: false },
                ...additionalScales.y, // Merge any additional y-axis options
            },
        },
        layout: {
            padding: {
                left: 10,
                right: 10,
                top: 10,
                bottom: 10,
            },
            ...restOptions.layout, // Merge any additional layout options
        },
        interaction: {
            mode: 'index',
            intersect: false,
            ...additionalOptions.interaction, // Merge any additional interaction options
        },
        datasets: {
            bar: {
                barPercentage: 0.6, // Adjusted for better visibility
                categoryPercentage: 0.7,
                ...additionalDatasets.bar, // Merge any additional bar dataset options
            },
            ...additionalDatasets, // Merge any additional datasets
        },
        ...restOptions, // Merge any other remaining options
    };
}

const predefinedColors = [
    '#FF6384', '#36A2EB', '#FFCE56', '#4BC0C0',
    '#9966FF', '#FF9F40', '#E7E9ED', '#76D7C4',
    '#C0392B', '#8E44AD', '#2ECC71', '#1ABC9C',
    '#3498DB', '#F1C40F', '#E67E22', '#95A5A6',
];

const shipColorMap = {};

/**
 * Assigns and retrieves a color for a given ship name.
 * @param {string} shipName - The name of the ship.
 * @returns {string} - The HEX color code.
 */
export function getShipColor(shipName) {
    if (shipColorMap[shipName]) {
        return shipColorMap[shipName];
    }
    const color = predefinedColors[Object.keys(shipColorMap).length % predefinedColors.length];
    shipColorMap[shipName] = color;
    return color;
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

