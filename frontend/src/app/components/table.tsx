import { useState, useEffect } from "react";

export const Table = (props: any) => {
    const [tableRows, setTableRows] = useState<JSX.Element[]>([]);
    
    useEffect(() => {
        const rows = props.data.map((row: any) => (
            <tr key={row.key} className="bg-white border-b dark:bg-gray-800 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-600">
                <th scope="row" className="px-6 py-4 font-medium text-gray-900 whitespace-nowrap dark:text-white">
                    {row.transaction_name.slice(0, 29)}
                </th>
                <td className="px-6 py-4">
                    {row.date}
                </td>
                <td className="px-6 py-4">
                    {row.category}
                </td>
                <td className={row.amount > 0 ? "px-6 py-4 text-red-400 font-bold" : "px-6 py-4 text-green-400 font-bold"}>
                    ${Math.abs(row.amount)}
                </td>
            </tr>
        ));
        setTableRows(rows.slice(0,5));
    }, [props.data]);
    
    return (
        <div className="relative overflow-x-auto shadow-md sm:rounded-lg">
            <table className="w-full text-sm text-left rtl:text-right text-gray-500 dark:text-gray-400">
                <thead className="text-xs text-gray-700 uppercase bg-gray-50 dark:bg-gray-700 dark:text-gray-400">
                    <tr>
                        <th scope="col" className="px-6 py-3">
                            Transaction name
                        </th>
                        <th scope="col" className="px-6 py-3">
                            Date
                        </th>
                        <th scope="col" className="px-6 py-3">
                            Category
                        </th>
                        <th scope="col" className="px-6 py-3">
                            Amount
                        </th>
                    </tr>
                </thead>
                <tbody>
                    {tableRows}
                </tbody>
            </table>
        </div>
    );
};
