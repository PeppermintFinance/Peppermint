export const Navbar = () => {
    return (
        <nav className="bg-white border-gray-200 dark:bg-gray-900">
        <div className="flex flex-wrap items-center justify-between mx-2 p-4">
            <a href="/" className="flex items-center space-x-3 rtl:space-x-reverse">
                <img src="https://img.icons8.com/plasticine/200/mint.png" className="h-8" alt="Peppermint Logo" />
                <span className="self-center text-2xl font-semibold whitespace-nowrap dark:text-white">Peppermint</span>
            </a>
            <div className="hidden w-full md:block md:w-auto" id="navbar-default">
                <ul className="flex flex-row">
                    <li>
                    <a href="/" className="block px-2 py-3 text-xl font-medium hover:text-green-600" aria-current="page">Home</a>
                    </li>
                    <li>
                    <a href="/insights" className="block px-2 py-3 text-xl font-medium hover:text-green-600 ">Insights</a>
                    </li>
                </ul>
            </div>
        </div>
        </nav>
    )
}