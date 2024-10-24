"use client";
import { Navbar } from "./components/navbar";
import { useState, useEffect } from "react";
import { usePlaidLink } from 'react-plaid-link';
import './globals.css'
export default function Home() {

  const [linkToken, setLinkToken] = useState(null)

  useEffect(() => {    
    const fetchLinkToken = async () => {
      const response = await fetch("http://localhost:8080/api/create-link-token", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({UserId: 'ladur'})
      });

      const data = await response.json();
      setLinkToken(data.linkToken);
    };

    fetchLinkToken();
  }, []); 

  const { open, ready } = usePlaidLink({
    token: linkToken,
    onSuccess: (public_token, metadata) => {
      fetch("http://localhost:8080/api/create-item", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({UserId: 'ladur', PublicToken: public_token, Metadata: metadata} )
      });
    },
  });

  const createUser = () => {
    fetch("http://localhost:8080/api/create-user", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({UserId: 'ladur', Username: 'ladur'})
    });
  }

  const refreshTransactions = () => {
    fetch("http://localhost:8080/api/refresh-user-items", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({UserId: 'ladur'})
    });
  }

  return (
    <div className="">
      <Navbar/>
      <section className="py-10 sm:py-16 lg:py-24">
        <div className="px-4 mx-auto max-w-7xl sm:px-6 lg:px-8">
            <div className="flex flex-row items-center gap-12">
                <div>
                    <h1 className="text-4xl font-bold text-black sm:text-6xl lg:text-8xl">
                        Embark on a financial odyssey.
                    </h1>
                    <p className="mt-8 text-base text-black sm:text-xl lg:text-2xl">No more wondering where your paycheck is going. Peppermint breaks it down for you so you can spend more time doing cool stuff, like building wooden horses.</p>
                    <div className="mt-10 sm:flex sm:items-center sm:space-x-8">
                    <button type="button" onClick={() => open()} disabled={!ready || !linkToken} className="w-52 text-white bg-gray-800 hover:bg-gray-900 focus:outline-none focus:ring-4 focus:ring-gray-300 font-medium rounded-full text-lg px-5 py-2.5 me-2 mb-2 dark:bg-gray-800 dark:hover:bg-gray-700 dark:focus:ring-gray-700 dark:border-gray-700">Add An Account</button>
                    <button type="button" onClick={createUser} className="w-52 text-white bg-gray-800 hover:bg-gray-900 focus:outline-none focus:ring-4 focus:ring-gray-300 font-medium rounded-full text-lg px-5 py-2.5 me-2 mb-2 dark:bg-gray-800 dark:hover:bg-gray-700 dark:focus:ring-gray-700 dark:border-gray-700">create user [test]</button>
                    <button type="button" onClick={refreshTransactions} className="w-52 text-white bg-gray-800 hover:bg-gray-900 focus:outline-none focus:ring-4 focus:ring-gray-300 font-medium rounded-full text-lg px-5 py-2.5 me-2 mb-2 dark:bg-gray-800 dark:hover:bg-gray-700 dark:focus:ring-gray-700 dark:border-gray-700">refresh transactions [test]</button>
                    </div>
                </div>
                <img src="cookie-bank.png" className="w-full h-auto" alt="Cookie Piggy Bank Image" />
            </div>
        </div>
    </section>
    </div>
  );
}