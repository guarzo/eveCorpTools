// static/js/chartConfigs/combinedLossesChartConfig.js
import { truncateLabel, getColor, getCommonOptions, validateChartData } from '../utils.js';

const combinedLossesChartConfig = {
    id: 'combinedLossesChart',
    instance: null,
    dataKeys: {
        mtd: 'mtdOurLossesValueData',
        ytd: 'ytdOurLossesValueData',
        lastMonth: 'lastMOurLossesValueData',
    },
    type: 'bar', // Change to 'pie' or 'doughnut' if preferred
    options: getCommonOptions('Combined Losses', {
        plugins: {
            legend: { display: true, position: 'top', labels: { color: '#ffffff' } },
            tooltip: {
                callbacks: {
                    label: function (context) {
                        const dataPoint = context.raw; // Access the data point directly
                        const lossesValue = dataPoint.LossesValue || 0;
                        const lossesCount = dataPoint.LossesCount || 0;
                        const shipType = dataPoint.ShipType || 'Unknown';
                        const shipCount = dataPoint.ShipCount || 0;
                        return `Value: ${lossesValue.toLocaleString()}, Count: ${lossesCount}, Ship: ${shipType} (${shipCount})`;
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
            },
            y: {
                beginAtZero: true,
                ticks: { color: '#ffffff' },
                grid: { color: '#444' },
                title: {
                    display: true,
                    text: 'Losses Value',
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

        // Ensure data has required fields
        // Expected fields: CharacterName, LossesValue, LossesCount, ShipType, ShipCount
        // If the backend doesn't provide all, adjust accordingly

        // Sort data by LossesValue descending
        const sortedData = [...data].sort((a, b) => (b.LossesValue || 0) - (a.LossesValue || 0));

        // Limit to top 10 characters
        const topN = 10;
        const limitedData = sortedData.slice(0, topN);

        const labels = limitedData.map(item => item.CharacterName || 'Unknown');
        const truncatedLabels = labels.map(label => truncateLabel(label, 15));

        const lossesValue = limitedData.map(item => item.LossesValue || 0);
        const lossesCount = limitedData.map(item => item.LossesCount || 0);
        const shipTypes = limitedData.map(item => item.ShipType || 'Unknown');
        const shipCounts = limitedData.map(item => item.ShipCount || 0);

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
        ];

        return { labels: truncatedLabels, datasets, fullLabels: labels };
    },
};

export default combinedLossesChartConfig;
