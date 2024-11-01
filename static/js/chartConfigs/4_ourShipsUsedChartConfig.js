// static/js/chartConfigs/4_ourShipsUsedChartConfig.js

import { truncateLabel, getShipColor, getCommonOptions, validateOurShipsUsedData } from '../utils.js';

/**
 * Configuration for the Our Ships Used Chart
 */
const ourShipsUsedChartConfig = {
    id: 'ourShipsUsedChart',
    instance: {}, // Correctly initialize as an object
    dataKeys: {
        mtd: { dataVar: 'mtdOurShipsUsedData', canvasId: 'ourShipsUsedChart_mtd' },
        ytd: { dataVar: 'ytdOurShipsUsedData', canvasId: 'ourShipsUsedChart_ytd' },
        lastMonth: { dataVar: 'lastMOurShipsUsedData', canvasId: 'ourShipsUsedChart_lastM' },
    },
    type: 'bar',
    dataType: 'object', // Specify that this chart expects object data
    options: getCommonOptions('Our Ships Used', {
        indexAxis: 'y',
        plugins: {
            tooltip: {
                mode: 'nearest', // Focus on the hovered bar segment
                intersect: true, // Show tooltip only when directly hovering over a segment
                callbacks: {
                    label: function (context) {
                        const value = context.parsed.x; // For horizontal bar chart
                        if (value > 0) { // Only show ships with count > 0
                            const shipName = context.dataset.label;
                            const index = context.dataIndex;
                            const total = context.chart.config.data.datasets.reduce((sum, dataset) => sum + (dataset.data[index] || 0), 0);
                            const percentage = total > 0 ? ((value / total) * 100).toFixed(2) : '0.00';
                            return `${shipName}: ${value} (${percentage}%)`;
                        } else {
                            return null; // Exclude ships with count <= 0 from tooltip
                        }
                    },
                },
            },
            datalabels: {
                color: '#ffffff',
                anchor: 'end',
                align: 'right',
                formatter: (value, context) => {
                    const index = context.dataIndex;
                    const total = context.chart.config.data.datasets.reduce((sum, dataset) => sum + (dataset.data[index] || 0), 0);
                    const percentage = total > 0 ? ((value / total) * 100).toFixed(1) : '0.0';
                    return `${value} (${percentage}%)`;
                },
                font: {
                    size: 10,
                    weight: 'bold',
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
        const chartName = 'Our Ships Used Chart';
        if (!validateOurShipsUsedData(data, chartName)) {
            // Trigger noData plugin
            return { labels: [], datasets: [] };
        }

        const characters = data.Characters || [];
        const shipNames = data.ShipNames || [];
        const seriesData = data.SeriesData || {};

        // Calculate total usage for each ship
        const shipUsage = shipNames.map(shipName => {
            const total = seriesData[shipName]?.reduce((a, b) => a + b, 0) || 0;
            return { shipName, total };
        });

        // Sort ships by total usage descending and limit to top 10
        const topShips = shipUsage
            .sort((a, b) => b.total - a.total)
            .slice(0, 10)
            .map(ship => ship.shipName);

        // Update shipNames to topShips
        const limitedShipNames = topShips;

        // Recalculate seriesData for limited ships
        const limitedSeriesData = {};
        limitedShipNames.forEach(shipName => {
            limitedSeriesData[shipName] = seriesData[shipName] || [];
        });

        const labels = characters.map(label => truncateLabel(label, 10));

        // Create datasets for each limited ship type
        const datasets = limitedShipNames.map((shipName) => ({
            label: shipName,
            data: limitedSeriesData[shipName] || [],
            backgroundColor: getShipColor(shipName),
            borderColor: '#ffffff',
            borderWidth: 1,
        }));

        return { labels, datasets };
    },
};

export default ourShipsUsedChartConfig;
