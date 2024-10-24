import Chart from "react-apexcharts";

export type DonutChartProps = {
    series: Array<Number>,
    labels: Array<string>,
  }

export const DonutChart = (props: DonutChartProps) => {
    const options = {
        series: props.series,
        colors: ["#1C64F2", "#16BDCA", "#FDBA8C", "#E74694"],
        chart: {
          height: 320,
          width: 400,
        },
        stroke: {
          colors: ["transparent"],
        },
        plotOptions: {
          pie: {
            donut: {
              labels: {
                show: true,
                name: {
                  show: true,
                  fontFamily: "Inter, sans-serif",
                  offsetY: 20,
                },
                total: {
                  showAlways: true,
                  show: true,
                  label: "Spending this month",
                  fontFamily: "Inter, sans-serif",
                  formatter: function (w: any) {
                    const sum = w.globals.seriesTotals.reduce((a: any, b: any) => {
                      return a + b
                    }, 0)
                    return '$' + sum.toFixed(2)
                  },
                },
                value: {
                  show: true,
                  fontFamily: "Inter, sans-serif",
                  offsetY: -20,
                  formatter: function (value: any) {
                    return '$' + parseFloat(value).toFixed(2)
                  },
                },
              },
              size: "80%",
            },
          },
        },
        grid: {
          padding: {
            top: -2,
          },
        },
        labels: props.labels,
        dataLabels: {
          enabled: false,
        },
        legend: {
          fontFamily: "Inter, sans-serif",
          position: "bottom",
        },
        yaxis: {
          labels: {
            formatter: function (value: any) {
              return '$' + parseFloat(value).toFixed(2)
            },
          },
        },
        xaxis: {
          labels: {
            formatter: function (value: any) {
              return '$' + parseFloat(value).toFixed(2)
            },
          },
          axisTicks: {
            show: false,
          },
          axisBorder: {
            show: false,
          },
        },
      }

    return (
        <div className="max-w-sm w-full bg-white rounded-lg shadow dark:bg-gray-800">
            <div className="py-6" id="donut-chart">
                <Chart options={options} height={options.chart.height} width={options.chart.width} series={options.series} type="donut"/>
            </div>
        </div>
    )
}