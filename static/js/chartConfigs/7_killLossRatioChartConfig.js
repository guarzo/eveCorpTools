// static/js/chartConfigs/7_killLossRatioChartConfig.js

import { truncateLabel, getColor, getCommonOptions, validateChartDataArray } from '../utils.js';

/**
 * Configuration for the Kill-to-Loss Ratio Chart
 */
const killLossRatioChartConfig = {
    id: 'killLossRatioChart',
    instance: {}, // Initialize as an object to store chart instances per timeframe
    dataKeys: {
        mtd: { dataVar: 'mtdKillLossRatioData', canvasId: 'killLossRatioChart_mtd' },
        ytd: { dataVar: 'ytdKillLossRatioData', canvasId: 'killLossRatioChart_ytd' },
        lastMonth: { dataVar: 'lastMKillLossRatioData', canvasId: 'killLossRatioChart_lastM' },
    },
    type: 'bar',
    dataType: 'array', // Specify that this chart expects array data
    options: getCommonOptions('Kill-to-Loss Ratio', {
        plugins: {
            legend: { display: false },
            tooltip: {
                callbacks: {
                    label: function (context) {
                        const index = context.dataIndex;
                        const kills = context.chart.config.data.additionalData.kills[index] || 0;
                        const losses = context.chart.config.data.additionalData.losses[index] || 0;
                        const ratio = context.parsed.y !== undefined ? context.parsed.y.toFixed(2) : '0.00';
                        return `Kills: ${kills}, Losses: ${losses}, Ratio: ${ratio}`;
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
                    text: 'Ratio',
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
        const chartName = 'Kill-to-Loss Ratio Chart';
        if (!validateChartDataArray(data, chartName)) {
            // Return empty labels and datasets to trigger the noDataPlugin
            return { labels: [], datasets: [], noDataMessage: 'No data available for this chart.' };
        }

        // Separate persistent characters and top ratios
        const persistentCharacters = data.filter(item => isPersistentCharacter(item.CharacterName));
        const topRatios = data.filter(item => !isPersistentCharacter(item.CharacterName))
            .sort((a, b) => (b.Ratio || 0) - (a.Ratio || 0))
            .slice(0, 10); // Top 10 by ratio

        // Combine persistent and top ratios
        const combinedData = [...persistentCharacters, ...topRatios];

        // Truncate if necessary
        const maxDisplay = 20;
        const limitedData = combinedData.slice(0, maxDisplay);

        // Check if there are at least 7 characters to display
        if (limitedData.length < 3) {
            console.warn(`Not enough data points (${limitedData.length}) for ${chartName}.`);
            return { labels: [], datasets: [], noDataMessage: 'Not enough data to display the chart.' };
        }

        const labels = limitedData.map(item => item.CharacterName || 'Unknown');
        const truncatedLabels = labels.map(label => truncateLabel(label, 15));

        const ratios = limitedData.map(item => item.Ratio || 0);
        const kills = limitedData.map(item => item.Kills || 0);
        const losses = limitedData.map(item => item.Losses || 0);

        const datasets = [{
            label: 'Kill-to-Loss Ratio',
            data: ratios,
            backgroundColor: ratios.map(ratio => getColor(ratio)),
            borderColor: 'rgba(75, 192, 192, 1)',
            borderWidth: 1,
        }];

        // Store kills and losses in a separate array for tooltip access
        const additionalData = {
            kills: kills,
            losses: losses,
        };

        return { labels: truncatedLabels, datasets, additionalData };
    },
};

// Helper function to determine if a character is persistent
function isPersistentCharacter(characterName) {
    // Define your logic to determine persistent characters
    // For example, based on a list of character names
    const persistentList = ['Gustav Oswaldo', 'Faoble', 'Stealthbot']; // Example list
    return persistentList.includes(characterName);
}

export default killLossRatioChartConfig;
