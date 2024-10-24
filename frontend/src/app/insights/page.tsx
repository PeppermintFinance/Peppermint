"use client"

import { AreaChart } from "../components/areachart";
import { Table } from "../components/table";
import { DonutChart } from "../components/donutchart";
import { useState, useEffect } from "react";
import { Navbar } from "../components/navbar";
import '../globals.css'

export default function Page () {
    const [tableData, setTableData] = useState([])
    const [netWorth, setNetWorth] = useState(0)
    const [spendingMap, setSpendingMap] = useState({})

    useEffect(() => {    
        const fetchTransactions = async () => {
          const response = await fetch('http://localhost:8080/api/get-transactions?' + new URLSearchParams({
            userId: 'ladur'
          }))
          const transactionData = await response.json()
    
          const tableData = transactionData.map((row: any, index: any) => ({
            transaction_name: row.Name,
            date: row.Date,
            category: row.Category,
            amount: row.Amount,
            key: index
          }))
    
          setTableData(tableData)
    
          // Filter by category
          const spendingMap: { [key: string]: number } = {};
          
          transactionData.forEach((transaction: any) => {
            // convert SCREAMING_SNAKE_CASE to Title Case
            const category = transaction.Category.toLowerCase().replace(/^_*(.)|_+(.)/g, (s: string, c: string, d: string) => c ? c.toUpperCase() : ' ' + d.toUpperCase())
            const amount = transaction.Amount;
            if (amount > 0) {
              spendingMap[category] = (spendingMap[category] || 0) + amount;
            }
          });
    
    
          setSpendingMap(spendingMap)
        };
    
    
        const fetchNetWorth = async () => {
          const response = await fetch('http://localhost:8080/api/get-net-worth?' + new URLSearchParams({
            userId: 'ladur'
          }))
          const data = await response.json()
          setNetWorth(data.NetWorth)
        };
    
        fetchTransactions();
        fetchNetWorth()
      }, []); 

    return (
      <div>
      <Navbar/>
      <div className="flex flex-row flex-wrap justify-between mx-10 my-10 p-5">
        <AreaChart data={[6500, 6418, 6456, 6526, 6356, 6456]} categories={['01 February', '02 February', '03 February', '04 February', '05 February', '06 February', '07 February']} percentage="12" color="#4ade80" headline={`$${netWorth.toFixed(2)}`} text="Net worth" />
        <Table data={tableData} />
        <div className="my-10">
          <DonutChart series={Object.values(spendingMap)} labels={Object.keys(spendingMap)}/>
        </div>
      </div>
      </div>
    )
}