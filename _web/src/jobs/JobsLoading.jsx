import React from 'react'
import LoadingIndicator from '../LoadingIndicator.jsx'

export const loadingJobs = () => (
  <section id="main-content-loading">
    <h2>Loading jobs ...</h2>
    <LoadingIndicator />
  </section>
)
