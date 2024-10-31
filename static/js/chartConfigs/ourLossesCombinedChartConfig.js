// static/js/chartConfigs/combinedLossesChartConfig.js

import { truncateLabel, getColor, getCommonOptions, validateChartData } from '../utils.js';

/**
 * Configuration for the Combined Losses Chart
 */
const combinedLossesChartConfig = {
    id: 'combinedLossesChart',
    instance: null,
    dataKeys: {
        mtd: 'mtdCombinedLossesData',
        ytd: 'ytdCombinedLossesData',
        lastMonth: 'lastMCombinedLossesData',
    },
    type: 'bar', // Adjust the chart type if needed (e.g., 'pie', 'bar')
    options: getCommonOptions('Combined Losses', {
        plugins: {
            legend: { display: true, position: 'top', labels: { color: '#ffffff' } },
            tooltip: {
                callbacks: {
                    label: function (context) {
                        const dataPoint = context.raw; // Access the data point directly
                        const lossesValue = dataPoint.lossesValue || 0;
                        const lossesCount = dataPoint.lossesCount || 0;
                        const shipType = dataPoint.shipType || 'Unknown';
                        const shipCount = dataPoint.shipCount || 0;
                        return `Value: ${lossesValue}, Count: ${lossesCount}, Ship: ${shipType} (${shipCount})`;
                    },
                },
            },
        },
        scales: {
            x: {
                type: 'category',
                ticks: {
                    color: '#ffffff',
                    maxRotation: 45,
                    minRotation: 45,
                    autoSkip: false,
                },
                grid: { display: false },
            },
            y: {
                beginAtZero: true,
                ticks: { color: '#ffffff' },
                grid: { color: '#444' },
                title: {
                    display: true,
                    text: 'Losses',
                    color: '#ffffff',
                    font: {
                        size: 14,
                        family: 'Montserrat, sans-serif',
                        weight: 'bold',
                    },
                },
            },
        },
    }),
    processData: function (data) {
        const chartName = 'Combined Losses Chart';
        if (!validateChartData(data, chartName)) {
            // Return empty labels and datasets to trigger the noDataPlugin
            return { labels: [], datasets: [] };
        }

        // Extract labels and data
        const labels = data.map(item => item.CharacterName || 'Unknown');
        const truncatedLabels = labels.map(label => truncateLabel(label, 10));

        const lossesValue = data.map(item => item.LossesValue || 0);
        const lossesCount = data.map(item => item.LossesCount || 0);
        const shipTypes = data.map(item => item.ShipType || 'Unknown');
        const shipCounts = data.map(item => item.ShipCount || 0);

        // Define datasets
        const datasets = [
            {
                label: 'Losses Value',
                data: lossesValue,
                backgroundColor: 'rgba(255, 99, 132, 0.7)',
                borderColor: 'rgba(255, 99, 132, 1)',
                borderWidth: 1,
            },
            {
                label: 'Losses Count',
                data: lossesCount,
                backgroundColor: 'rgba(54, 162, 235, 0.7)',
                borderColor: 'rgba(54, 162, 235, 1)',
                borderWidth: 1,
            },
            // Optionally, add a dataset for Ship Counts or use stacked bars
        ];

        return { labels: truncatedLabels, datasets, fullLabels: labels };
    },
};

export default combinedLossesChartConfig;
