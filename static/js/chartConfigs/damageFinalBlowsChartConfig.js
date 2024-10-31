// static/js/chartConfigs/damageFinalBlowsChartConfig.js

import { truncateLabel, getColor, getCommonOptions, validateChartData } from '../utils.js';

/**
 * Configuration for the Damage Done and Final Blows Chart
 */
const damageFinalBlowsChartConfig = {
    id: 'damageFinalBlowsChart',
    instance: null,
    dataKeys: {
        mtd: 'mtdCharacterDamageData',
        ytd: 'ytdCharacterDamageData',
        lastMonth: 'lastMCharacterDamageData',
    },
    type: 'bar', // Grouped bar chart
    options: getCommonOptions('Damage Done and Final Blows', {
        plugins: {
            legend: { display: true, position: 'top', labels: { color: '#ffffff' } },
            tooltip: {
                callbacks: {
                    label: function (context) {
                        const label = context.dataset.label || '';
                        const value = context.parsed.y !== undefined ? context.parsed.y : context.parsed.x;
                        return `${label}: ${value}`;
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
                    text: 'Count',
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
        const chartName = 'Damage Done and Final Blows';
        if (!validateChartData(data, chartName)) {
            // Return empty labels and datasets to trigger the noDataPlugin
            return { labels: [], datasets: [] };
        }

        // Extract labels and data
        const labels = data.map(item => item.Name || 'Unknown');
        const truncatedLabels = labels.map(label => truncateLabel(label, 10));

        const damageData = data.map(item => item.DamageDone || 0);
        const finalBlowsData = data.map(item => item.FinalBlows || 0);

        // Define datasets
        const datasets = [
            {
                label: 'Damage Done',
                data: damageData,
                backgroundColor: 'rgba(255, 77, 77, 0.7)',
                borderColor: 'rgba(255, 77, 77, 1)',
                borderWidth: 1,
            },
            {
                label: 'Final Blows',
                data: finalBlowsData,
                backgroundColor: 'rgba(54, 162, 235, 0.7)',
                borderColor: 'rgba(54, 162, 235, 1)',
                borderWidth: 1,
            },
        ];

        return { labels: truncatedLabels, datasets, fullLabels: labels };
    },
};

export default damageFinalBlowsChartConfig;
