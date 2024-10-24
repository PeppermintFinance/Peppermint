export const metadata = {
  title: 'Peppermint',
  description: 'Personal Finance Aggregator',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en" className="">
      <body>{children}</body>
    </html>
  )
}
