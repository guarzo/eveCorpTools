// chartConfigs/ourShipsUsedChartConfig.js
import { truncateLabel, getColor, commonOptions } from '../utils.js';

const ourShipsUsedChartConfig = {
    id: 'ourShipsUsedChart',
    instance: null,
    dataKeys: {
        mtd: 'mtdOurShipsUsedData',
        ytd: 'ytdOurShipsUsedData',
        lastMonth: 'lastMOurShipsUsedData',
    },
    type: 'bar',
    options: {
        ...commonOptions,
        indexAxis: 'y',
        plugins: {
            ...commonOptions.plugins,
            tooltip: {
                ...commonOptions.plugins.tooltip,
                callbacks: {
                    ...commonOptions.plugins.tooltip.callbacks,
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
                ...commonOptions.scales.x,
                stacked: true,
            },
            y: {
                ...commonOptions.scales.y,
                stacked: true,
                ticks: {
                    ...commonOptions.scales.y.ticks,
                    autoSkip: false,
                },
                grid: { display: false },
            },
        },
    },
    processData: function (data) {
        // ... processing logic remains the same ...
        const characters = data.Characters || [];
        const shipNames = data.ShipNames || [];
        const seriesData = data.SeriesData || {};

        const fullLabels = [...characters];
        const labels = characters.map(label => truncateLabel(label, 10));

        const datasets = shipNames.map((shipName, index) => ({
            label: shipName,
            data: seriesData[shipName] || [],
            backgroundColor: getColor(index),
        }));

        return { labels, datasets, fullLabels };
    },
};

export default ourShipsUsedChartConfig;
