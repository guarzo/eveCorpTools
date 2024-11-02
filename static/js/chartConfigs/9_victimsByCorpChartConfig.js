// static/js/chartConfigs/victimsByCorpChartConfig.js

import { truncateLabel, getColor, getCommonOptions, validateChartDataArray } from '../utils.js';

/**
 * Configuration for the Kills by Corporation Bar Chart
 */
const victimsByCorpChartConfig = {
    id: 'victimsByCorpChart',
    instance: {}, // To store chart instances per timeframe if needed
    dataKeys: {
        mtd: { dataVar: 'mtdVictimsByCorpData', canvasId: 'victimsByCorpChart_mtd' },
        ytd: { dataVar: 'ytdVictimsByCorpData', canvasId: 'victimsByCorpChart_ytd' },
        lastMonth: { dataVar: 'lastMVictimsByCorpData', canvasId: 'victimsByCorpChart_lastM' },
    },
    type: 'bar',
    dataType: 'array', // Expecting an array of CorporationKillCount
    options: getCommonOptions('Victims by Corporation', {
        plugins: {
            legend: { display: false },
            tooltip: {
                callbacks: {
                    label: function (context) {
                        const killCount = context.parsed.y !== undefined ? context.parsed.y : 0;
                        return `Killmails: ${killCount}`;
                    },
                },
            },
        },
        scales: {
            x: {
                type: 'category',
                title: {
                    display: true,
                    text: 'Corporations',
                    color: '#ffffff',
                    font: {
                        size: 14,
                        family: 'Montserrat, sans-serif',
                        weight: 'bold',
                    },
                },
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
                title: {
                    display: true,
                    text: 'Kill Count',
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
        },
    }),
    processData: function (data) {
        const chartName = 'Victims by Corporation Chart';
        if (!validateChartDataArray(data, chartName)) {
            console.log("failed validate for Victims by Corp")
            console.log(data,chartName)
            // Return empty labels and datasets to trigger the noDataPlugin
            return { labels: [], datasets: [] };
        }

        // Sort data by kill count descending
        const sortedData = data.sort((a, b) => b.KillCount - a.KillCount);

        // Limit to top 15 corporations
        const limitedData = sortedData.slice(0, 15);

        // Prepare labels and data
        const labels = limitedData.map(item => item.name || item.corporation_id || "Unknown");
        const counts = limitedData.map(item => item.kill_count);

        // Assign colors to each bar
        const backgroundColors = counts.map((count, index) => getColor(index));

        const datasets = [{
            label: 'Killmails',
            data: counts,
            backgroundColor: backgroundColors,
            borderColor: 'rgba(75, 192, 192, 1)',
            borderWidth: 1,
        }];
        console.log("processed data victims by corp", labels, datasets)

        return { labels, datasets };
    },
};

export default victimsByCorpChartConfig;
