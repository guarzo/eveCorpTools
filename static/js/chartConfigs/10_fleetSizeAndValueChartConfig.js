// static/js/chartConfigs/fleetSizeAndValueChartConfig.js

import { truncateLabel, getColor, getCommonOptions, validateChartDataArray } from '../utils.js';

/**
 * Configuration for the Average Fleet Size and Total Value Over Time Chart
 */
const fleetSizeAndValueChartConfig = {
    id: 'fleetSizeAndValueChart',
    instance: {}, // Initialize as an object to store chart instances per timeframe
    dataKeys: {
        mtd: { dataVar: 'mtdFleetSizeAndValueData', canvasId: 'fleetSizeAndValueChart_mtd' },
        ytd: { dataVar: 'ytdFleetSizeAndValueData', canvasId: 'fleetSizeAndValueChart_ytd' },
        lastMonth: { dataVar: 'lastMFleetSizeAndValueData', canvasId: 'fleetSizeAndValueChart_lastM' },
    },
    type: 'line', // Using line chart
    dataType: 'array', // Updated to 'array' since backend returns a single array
    options: getCommonOptions('Fleet Size and Value Killed Over Time', {
        plugins: {
            legend: { display: true },
            tooltip: {
                mode: 'index',
                intersect: false,
                callbacks: {
                    label: function (context) {
                        const datasetLabel = context.dataset.label || '';
                        const value = context.parsed.y !== undefined ? context.parsed.y : 0;
                        return `${datasetLabel}: ${value}`;
                    },
                },
            },
        },
        scales: {
            x: {
                type: 'time',
                time: {
                    unit: 'day', // Change to 'week' if interval is weekly
                    tooltipFormat: 'MMM dd, yyyy',
                },
                title: {
                    display: true,
                    text: 'Time',
                    color: '#ffffff',
                    font: {
                        size: 14,
                        family: 'Montserrat, sans-serif',
                        weight: 'bold',
                    },
                },
                ticks: {
                    color: '#ffffff',
                },
                grid: { display: false },
            },
            y: {
                type: 'linear',
                position: 'left',
                beginAtZero: true,
                title: {
                    display: true,
                    text: 'Average Fleet Size',
                    color: '#ffffff',
                    font: {
                        size: 14,
                        family: 'Montserrat, sans-serif',
                        weight: 'bold',
                    },
                },
                ticks: {
                    color: '#ffffff',
                },
                grid: { color: '#444' },
            },
            y1: {
                type: 'linear',
                position: 'right',
                beginAtZero: true,
                title: {
                    display: true,
                    text: 'Total Value',
                    color: '#ffffff',
                    font: {
                        size: 14,
                        family: 'Montserrat, sans-serif',
                        weight: 'bold',
                    },
                },
                ticks: {
                    color: '#ffffff',
                },
                grid: { display: false },
            },
        },
    }),
    processData: function (data) {
        const chartName = 'Average Fleet Size and Total Value Over Time Chart';
        if (!validateChartDataArray(data, chartName)) {
            // Return empty labels and datasets to trigger the noDataPlugin
            return { labels: [], datasets: [] };
        }

        // Prepare labels and data
        const labels = data.map(item => new Date(item.time));

        const fleetSizes = data.map(item => item.avg_fleet_size || 0);
        const totalValues = data.map(item => item.total_value || 0);

        // Check if there are at least 7 characters to display
        if (fleetSizes.length < 3) {
            console.warn(`Not enough data points (${fleetSizes.length}) for ${chartName}.`);
            return { labels: [], datasets: [], noDataMessage: 'Not enough data to display the chart.' };
        }

        const datasets = [
            {
                label: 'Average Fleet Size',
                data: fleetSizes,
                backgroundColor: 'rgba(75, 192, 192, 0.7)',
                borderColor: 'rgba(75, 192, 192, 1)',
                borderWidth: 2,
                fill: true,
                tension: 0.1,
                yAxisID: 'y',
            },
            {
                label: 'Total Value',
                data: totalValues,
                backgroundColor: 'rgba(255, 159, 64, 0.7)',
                borderColor: 'rgba(255, 159, 64, 1)',
                borderWidth: 2,
                fill: false,
                tension: 0.1,
                yAxisID: 'y1',
            }
        ];

        return { labels, datasets };
    },
};

export default fleetSizeAndValueChartConfig;
