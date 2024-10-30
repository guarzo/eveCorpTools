// static/js/chartConfigs/ourShipsUsedChartConfig.js

import { truncateLabel, getColor, getCommonOptions } from '../utils.js';

/**
 * Configuration for the Our Ships Used Chart
 */
const ourShipsUsedChartConfig = {
    id: 'ourShipsUsedChart',
    instance: null,
    dataKeys: {
        mtd: 'mtdOurShipsUsedData',
        ytd: 'ytdOurShipsUsedData',
        lastMonth: 'lastMOurShipsUsedData',
    },
    type: 'bar',
    options: getCommonOptions('Our Ships Used', {
        indexAxis: 'y',
        plugins: {
            tooltip: {
                callbacks: {
                    label: function (context) {
                        const shipName = context.dataset.label;
                        const value = context.parsed.x;
                        return `${shipName}: ${value}`;
                    },
                },
            },
        },
        scales: {
            x: {
                stacked: true,
                ticks: { color: '#ffffff' },
                grid: { display: false },
            },
            y: {
                stacked: true,
                ticks: {
                    color: '#ffffff',
                    autoSkip: false,
                },
                grid: { display: false },
            },
        },
    }),
    processData: function (data) {
        // Extract characters, ship names, and series data
        const characters = data.Characters || [];
        const shipNames = data.ShipNames || [];
        const seriesData = data.SeriesData || {};

        const fullLabels = [...characters];
        const labels = characters.map(label => truncateLabel(label, 10));

        // Create datasets for each ship type
        const datasets = shipNames.map((shipName, index) => ({
            label: shipName,
            data: seriesData[shipName] || [],
            backgroundColor: getColor(index),
        }));

        return { labels, datasets, fullLabels };
    },
};

export default ourShipsUsedChartConfig;
