import React, { useState } from 'react'

import './App.css'
import Header from './Header'
import AuthenticatedPage from './Authentication'

export default function App() {
  const [isAuthenticated, setIsAuthenticated] = useState(false);

  return (
    <>
      <Header />
      <main>
        {!isAuthenticated ? (
          <AuthenticatedPage setIsAuthenticated={setIsAuthenticated} /> 
        ): 'user profile here'}
      </main>
    </>
  )
}
