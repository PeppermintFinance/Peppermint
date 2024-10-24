import React, { Component } from "react";
import Chart from "react-apexcharts";

// graph options from https://flowbite.com/docs/plugins/charts/


export type AreaChartProps = {
  color: string,
  percentage: string,
  headline: string,
  text: string,
  data: Array<Number>,
  categories: Array<String>,
}

export const AreaChart = (props: AreaChartProps) => {

  const options = {
    chart: {
      height: "100%",
      maxWidth: "100%",
      fontFamily: "Inter, sans-serif",
      dropShadow: {
        enabled: false,
      },
      toolbar: {
        show: false,
      },
    },
    tooltip: {
      enabled: true,
      x: {
        show: false,
      },
    },
    fill: {
      type: "gradient",
      gradient: {
        opacityFrom: 0.55,
        opacityTo: 0,
        shade: "#1C64F2",
        gradientToColors: ["#1C64F2"],
      },
    },
    dataLabels: {
      enabled: false,
    },
    stroke: {
      width: 6,
    },
    grid: {
      show: false,
      strokeDashArray: 4,
      padding: {
        left: 2,
        right: 2,
        top: 0
      },
    },
    series: [
      {
        name: "Net Worth",
        data: props.data,
        color: props.color,
      },
    ],
    xaxis: {
      categories: props.categories,
      labels: {
        show: false,
      },
      axisBorder: {
        show: false,
      },
      axisTicks: {
        show: false,
      },
    },
    yaxis: {
      show: false,
    },
  }

    return (
    <div className="max-w-sm w-full bg-white rounded-lg shadow dark:bg-gray-800 p-4 md:p-6">
      <div className="flex justify-between">
        <div>
          <h5 className="leading-none text-3xl font-bold text-gray-900 dark:text-white pb-2">{props.headline}</h5>
          <p className="text-base font-normal text-gray-500 dark:text-gray-400">{props.text}</p>
        </div>
        <div
          className="flex items-center px-2.5 py-0.5 text-base font-semibold text-blue-500 dark:text-blue-500 text-center">
          {props.percentage}%
          <svg className="w-3 h-3 ms-1" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 10 14">
            <path stroke="currentColor" strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M5 13V1m0 0L1 5m4-4 4 4"/>
          </svg>
        </div>
      </div>
      <Chart options={options} series={options.series} type="area"/>
    </div>
    );
}