const width = 960;
const height = 500;
const margin = 5;
const padding = 5;
const adj = 30;

const dayString = day => day <= 9 ? "0" + day : "" + day;

function convertDateString(dateString) {
  const newDate = new Date(dateString);
  return newDate.getFullYear() + "-" + (newDate.getMonth() + 1) + "-" + newDate.getDay();
}
function convertDate(date) {
  return date.getFullYear() + "-" + (date.getMonth() + 1) + "-" + date.getDay();
}

function getDataSet(options) {
  const columnsNames = options.map(o => o.OptionID);
  columnsNames.unshift("date")

  const dates = new Object();
  options.forEach(option => {
    //add fake data
    for (let i = 0; i <= 31; i += 1) {
      const fake_date = "2022-01-" + dayString(i);
      if (dates[fake_date] === undefined) dates[fake_date] = new Object();
      dates[fake_date][option.OptionID] = option.Prices[0].Price + Math.random() * 20 - 10;
    }
    option.Prices.forEach(
      p => {
        const optionID = option.OptionID;
        const date = convertDate(new Date(p.CheckDate));
        const price = p.Price;
        if (dates[date] === undefined) dates[date] = new Object();
        dates[date][optionID] = price;

      });
  });

  const rows = []
  Object.keys(dates).sort().forEach(d => {
    const row = new Object();

    Object.keys(dates[d]).forEach(id => {
      row[id] = dates[d][id];
    })

    row.date = d;
    rows.push(row);
  })
  rows.columns = columnsNames;
  return rows;
}
function parseData(raw_data) {
  return getDataSet(raw_data);
}

// we are appending SVG first
const svg = d3.select("div#container").append("svg")
  .attr("preserveAspectRatio", "xMinYMin meet")
  .attr("viewBox", "-"
    + adj + " -"
    + adj + " "
    + (width + adj * 3) + " "
    + (height + adj * 3))
  .style("padding", padding)
  .style("margin", margin)
  .classed("svg-content", true);

//-----------------------------DATA-----------------------------//
const timeConv = d3.timeParse("%Y-%m-%d");
const dataset = d3.csv("/static/data.csv");
fetch("/obtenerOpciones?ProductID=" + optionID)
  .then(response => response.json())
  .then(function(raw_data) {
    const data = parseData(raw_data);
    // dataset.then(function(data) {
    var slices = data.columns.slice(1).map(function(id) {
      return {
        id: id,
        values: data.map(function(d) {
          return {
            date: timeConv(d.date),
            measurement: +d[id]
          };
        })
      };
    });






    //----------------------------SCALES----------------------------//
    const xScale = d3.scaleTime().range([0, width]);
    const yScale = d3.scaleLinear().rangeRound([height, 0]);
    xScale.domain(d3.extent(data, function(d) {
      return timeConv(d.date)
    }));
    yScale.domain([(0), d3.max(slices, function(c) {
      return d3.max(c.values, function(d) {
        return d.measurement + 4;
      });
    })
    ]);

    //-----------------------------AXES-----------------------------//
    const yaxis = d3.axisLeft()
      .ticks((slices[0].values).length)
      .scale(yScale);

    const xaxis = d3.axisBottom()
      .ticks(d3.timeDay.every(1))
      .tickFormat(d3.timeFormat('%b %d'))
      .scale(xScale);

    //----------------------------LINES-----------------------------//
    const line = d3.line()
      .x(function(d) { return xScale(d.date); })
      .y(function(d) { return yScale(d.measurement); });

    let id = 0;
    const ids = function() {
      return "line line-" + id++;
    }
    //-------------------------2. DRAWING---------------------------//
    //-----------------------------AXES-----------------------------//
    svg.append("g")
      .attr("class", "axis")
      .attr("transform", "translate(0," + height + ")")
      .call(xaxis);

    svg.append("g")
      .attr("class", "axis")
      .call(yaxis)
      .append("text")
      .attr("transform", "rotate(-90)")
      .attr("dy", ".75em")
      .attr("y", 6)
      .style("text-anchor", "end")
      .text("Price");

    //----------------------------LINES-----------------------------//
    const lines = svg.selectAll("lines")
      .data(slices)
      .enter()
      .append("g");

    lines.append("path")
      .attr("class", ids)
      .attr("d", function(d) { return line(d.values); });

    lines.append("text")
      .attr("class", "serie_label")
      .datum(function(d) {
        return {
          id: d.id,
          value: d.values[d.values.length - 1]
        };
      })
      .attr("transform", function(d) {
        return "translate(" + (xScale(d.value.date) + 10)
          + "," + (yScale(d.value.measurement) + 5) + ")";
      })
      .attr("x", 5)
      .text(function(d) { return ("") + d.id; });

  });
