const cornerRadius = 1;

function HBar(id, width, min, max, data, value, label, click) {
    const barHeight = 18;
    const barPadding = 4;
    const barSeparation = barHeight + barPadding;

    var tooltip = d3.select('#tooltip')
        .style('opacity', 0);

    const scale = d3.scaleLinear()
        .domain([min, max])
        .range([0, width]);

    const chart = d3.select(id)
        .attr('width', width)
        .attr('height', barSeparation * data.length);

    chart.selectAll('g').remove();

    const bar = chart.selectAll('g')
        .data(data)
        .enter()
        .append('g')
        .attr('transform', function(data, index) {
            return 'translate(0,' + index * barSeparation + ')';
        })
        .on('mouseover', function(event, data) {
            tooltip.transition()
                .duration(50)
                .style('opacity', 1);
            tooltip.html(label(data) + '<br />' + value(data));
            positionTooltip(tooltip, event);
        })
        .on('mouseout', function(event, data) {
            tooltip.transition()
               .duration(50)
               .style('opacity', 0);
        })
        .on('click', click);

    bar.append('rect')
        .attr('rx', cornerRadius)
        .attr('width', function(data) {
            return scale(value(data));
        })
        .attr('height', barHeight);

    bar.append('text')
        .attr('x', barPadding)
        .attr('y', barHeight / 2)
        .attr('dy', '.35em')
        .attr('text-anchor', 'start')
        .text(label);

    bar.append('text')
        .attr('x', width-barPadding)
        .attr('y', barHeight / 2)
        .attr('dy', '.35em')
        .attr('text-anchor', 'end')
        .text(value);
}

function Histogram(id, width, height, min, max, data, value, label, click) {
    const margin = {top: 10, right: 10, bottom: 60, left: 60};

    var tooltip = d3.select('#tooltip')
        .style('opacity', 0);

    const x = d3.scaleUtc()
        .domain([new Date(min * 1000), new Date(max * 1000)])
        .range([0, width - margin.right - margin.left]);

    const chart = d3.select(id)
        .attr('width', width)
        .attr('height', height);

    chart.selectAll('g').remove();

    chart.append('g')
        .attr('transform', `translate(${margin.left}, ${height - margin.bottom})`)
        .call(d3.axisBottom(x))
        .selectAll('text')
        .style('text-anchor', 'end')
        .attr('dx', '-.8em')
        .attr('dy', '.15em')
        .attr('transform', 'rotate(-60)');

    const histogram = d3.histogram()
        .value(function(data) {
            return new Date(value(data) * 1000);
        })
        .domain(x.domain())
        .thresholds(x.ticks(Math.floor((width - margin.right - margin.left) / 9)));

    const bins = histogram(data);

    const y = d3.scaleLinear()
        .domain([0, d3.max(bins, function(data) {
            return data.length;
        })])
        .range([height - margin.top - margin.bottom, 0]);

    chart.append('g')
        .attr('transform', `translate(${margin.left}, ${margin.top})`)
        .call(d3.axisLeft(y));

    chart.append('g')
        .attr('transform', `translate(${margin.left}, ${margin.top})`)
        .selectAll('rect')
        .data(bins)
        .enter()
        .append('rect')
        .attr('rx', cornerRadius)
        .attr('x', 1)
        .attr('transform', function(data) {
            return `translate(${x(data.x0)}, ${y(data.length)})`;
        })
        .attr('width', function(data) {
            var w = x(data.x1) - x(data.x0) - 1;
            if (w < 0) {
                w = 0;
            }
            return w;
        })
        .attr('height', function(data) {
            return height - margin.top - margin.bottom - y(data.length);
        })
        .on('mouseover', function(event, data) {
            tooltip.transition()
                .duration(50)
                .style('opacity', 1);
            tooltip.html(label(data));
            positionTooltip(tooltip, event);
        })
        .on('mouseout', function(event, data) {
            tooltip.transition()
               .duration(50)
               .style('opacity', 0);
        })
        .on('click', click);
}

function positionTooltip(tooltip, event) {
    if (event.pageX < (window.innerWidth / 2)) {
        tooltip.style('left', event.pageX + 'px');
        tooltip.style('right', 'auto');
    } else {
        tooltip.style('left', 'auto');
        tooltip.style('right', (window.innerWidth - event.pageX) + 'px');
    }
    if (event.pageY < (window.innerHeight / 2)) {
        tooltip.style('top', event.pageY + 'px');
        tooltip.style('bottom', 'auto');
    } else {
        tooltip.style('top', 'auto');
        tooltip.style('bottom', (window.innerHeight - event.pageY) + 'px');
    }
}