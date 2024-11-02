// static/js/chartConfigs/4_ourShipsUsedChartConfig.js

import { truncateLabel, getShipColor, getCommonOptions, validateOurShipsUsedData } from '../utils.js';

/**
 * Configuration for the Our Ships Used Chart
 */
const ourShipsUsedChartConfig = {
    id: 'ourShipsUsedChart',
    instance: {}, // Initialize as an object to store chart instances per timeframe
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
                        const value = context.parsed.x !== undefined ? context.parsed.x : context.parsed.y;
                        const shipName = context.dataset.label || '';
                        return `${shipName}: ${value} Kills`;
                    },
                },
            },
            datalabels: {
                color: '#ffffff',
                anchor: 'end',
                align: 'right',
                formatter: (value) => `${value}`,
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
                title: {
                    display: true,
                    text: 'Kills',
                    color: '#ffffff',
                    font: {
                        size: 14,
                        family: 'Montserrat, sans-serif',
                        weight: 'bold',
                    },
                },
            },
            y: {
                stacked: true,
                ticks: {
                    color: '#ffffff',
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
        },
        responsive: true,
        maintainAspectRatio: false, // Allow the chart to adjust its height
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

        // Constants to define the maximum number of ships and characters to display
        const MAX_SHIPS = 10; // Adjust based on requirements
        const MAX_CHARACTERS = 15; // Adjust based on requirements

        // Calculate total usage for each ship
        const shipUsage = shipNames.map(shipName => {
            const total = seriesData[shipName]?.reduce((a, b) => a + b, 0) || 0;
            return { shipName, total };
        });

        // Sort ships by total usage descending and limit to top MAX_SHIPS
        const topShips = shipUsage
            .sort((a, b) => b.total - a.total)
            .slice(0, MAX_SHIPS)
            .map(ship => ship.shipName);

        // Update shipNames to topShips
        const limitedShipNames = topShips;

        // Recalculate seriesData for limited ships
        const limitedSeriesData = {};
        limitedShipNames.forEach(shipName => {
            limitedSeriesData[shipName] = seriesData[shipName] || [];
        });

        // Calculate total usage per character across all limited ships
        const characterUsage = characters.map((char, index) => {
            let total = 0;
            limitedShipNames.forEach(ship => {
                total += seriesData[ship]?.[index] || 0;
            });
            return { character: char, total };
        });

        // Sort characters by total usage descending and limit to top MAX_CHARACTERS
        const topCharacters = characterUsage
            .sort((a, b) => b.total - a.total)
            .slice(0, MAX_CHARACTERS)
            .map(item => item.character);

        // Find the indices of topCharacters in the original characters array
        const topCharacterIndices = topCharacters.map(char => characters.indexOf(char)).filter(index => index !== -1);

        // Prepare labels: truncate and limit to topCharacters
        const labels = topCharacters.map(label => truncateLabel(label, 10));

        // Create datasets for each limited ship type, only for topCharacters
        const datasets = limitedShipNames.map(shipName => ({
            label: shipName,
            data: topCharacterIndices.map(index => seriesData[shipName]?.[index] || 0),
            backgroundColor: getShipColor(shipName),
            borderColor: '#ffffff',
            borderWidth: 1,
        }));

        return { labels, datasets };
    },
};

export default ourShipsUsedChartConfig;
