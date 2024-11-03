// static/js/chartConfigs/topShipsKilledChartConfig.js

import { truncateLabel, getShipColor, getCommonOptions, validateChartDataArray } from '../utils.js';

/**
 * Configuration for the Top Ships Killed Chart
 */
const topShipsKilledChartConfig = {
    type: 'wordCloud',
    options: getCommonOptions('Top Ships Killed', {
        plugins: {
            legend: { display: false },
            tooltip: {
                callbacks: {
                    label: function (context) {
                        const word = context.label || 'Unknown';
                        const count = context.raw / context.scalingFactor || 0;
                        return `${word}: ${count} Killmails`;
                    },
                },
                mode: 'nearest',
                intersect: true,
                backgroundColor: '#000000',
                titleColor: '#ffffff',
                bodyColor: '#ffffff',
            },
        },
        scales: {
            x: { display: false },
            y: { display: false },
        },
        layout: { padding: 10 },
        animation: { duration: 0 },
        font: { weight: 'normal', family: 'Arial' },
    }),
    processData: function (data) {
        const chartName = 'Top Ships Killed Chart';
        if (!validateChartDataArray(data, chartName)) {
            // Return empty data to trigger the noDataPlugin
            return { labels: [], datasets: [] };
        }

        console.log('Incoming data for Top Ships Killed:', data); // Debugging log

        // Map data to {text: ShipName, weight: KillCount} format
        const mappedData = data.map(item => {
            if (!item || !item.Name || typeof item.KillCount !== 'number') {
                console.warn(`Invalid data point in "${chartName}":`, item);
                return null;
            }
            return {
                text: item.Name,
                weight: item.KillCount
            };
        }).filter(item => item !== null);

        if (mappedData.length === 0) {
            console.warn(`No valid data points for Word Cloud (${chartName}).`);
            return { labels: [], datasets: [] };
        }

        // Limit to top 10 ships
        const maxWords = 10;
        const limitedData = mappedData.slice(0, maxWords);

        // Find maximum kill count for scaling
        const maxKillCount = Math.max(...limitedData.map(d => d.weight));
        const scalingFactor = maxKillCount > 0 ? 60 / maxKillCount : 10;

        const scaledData = limitedData.map(d => ({
            text: d.text,
            weight: d.weight * scalingFactor
        }));

        const labels = scaledData.map(d => d.text);
        const weights = scaledData.map(d => d.weight);
        const colors = scaledData.map((d, index) => getShipColor(d.text, index));
        const rotations = scaledData.map(() => (Math.random() > 0.5 ? 0 : 90));

        const datasets = [{
            data: weights,
            color: colors,
            rotation: rotations,
            scalingFactor: scalingFactor, // Pass scalingFactor for tooltip calculations
        }];

        return { labels, datasets };
    },
};

export default topShipsKilledChartConfig;
