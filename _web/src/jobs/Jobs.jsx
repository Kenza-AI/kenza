import React from 'react'

import Error from '../Error.jsx'
import { jobsSection } from './JobsSection.jsx'
import { loadingJobs } from './JobsLoading.jsx'

export default class Jobs extends React.Component {

  // Component lifecycle
  componentDidMount() {
    this.fetchJobs()
    this.startPollingJobs()
    this.listenForKeyboadEvents()
  }

  componentWillUnmount() {
    this.stopPollingJobs()
    this.stopListeningForKeyboardEvents()
  }

  componentDidUpdate(prevProps) {
    const isDifferentProject = this.props.projectID != prevProps.projectID
    if (isDifferentProject) {
      this.props.fetchJobsHandler(this.props.projectID)
    }
  }

  // Component rendering
  render() {
    const {
      jobs,
      filterTextChangeHandler,
      viewJobDetailsHandler,
      viewLogsHandler,
      filtersChangeHandler,
      selectJobForComparisonHandler,
      compareJobsHandler,
      dismissJobsComparisonHandler,
      clearJobsSelectedForComparisonHandler,
      deleteJobsHandler,
      cancelJobsHandler } = this.props

    if (jobs.fetchJobsError != null) {
      return <Error message={jobs.fetchJobsError} />
    } else if (jobs.deleteJobsError != null) {
      return <Error message={jobs.deleteJobsError} />
    } else if (jobs.cancelJobsError != null) {
      return <Error message={jobs.cancelJobsError} />
    }  else if (jobs.isFetching && this.jobPollingIntervalID == undefined || jobs.isDeleting || jobs.isCancelling) {
      return loadingJobs()
    } else {
      return jobsSection(
        jobs,
        filterTextChangeHandler,
        viewJobDetailsHandler,
        viewLogsHandler,
        filtersChangeHandler,
        selectJobForComparisonHandler,
        compareJobsHandler,
        dismissJobsComparisonHandler,
        deleteJobsHandler,
        cancelJobsHandler)
    }
  }

  fetchJobs() {
    this.props.fetchJobsHandler(this.props.projectID)
  }

  listenForKeyboadEvents() {
    this.handleKeyboardEvent = this.handleKeyboardEvent.bind(this)
    document.addEventListener('keyup', this.handleKeyboardEvent) // one can clear jobs selected for comparison using the `Esc` key.
  }

  stopListeningForKeyboardEvents() {
    document.removeEventListener('keyup', this.handleKeyboardEvent)
  }

  handleKeyboardEvent(event) {
    const handler = this.props.jobs.inComparisonMode ?
      this.props.dismissJobsComparisonHandler : this.props.clearJobsSelectedForComparisonHandler
    event.key === "Escape" && handler()
  }

  startPollingJobs() {
    const skipFetch = jobs => (jobs.inComparisonMode)
    this.jobPollingIntervalID = setInterval(
      () => skipFetch(this.props.jobs) || this.fetchJobs(), 15000) // 15 seconds
  }

  stopPollingJobs() {
    clearInterval(this.jobPollingIntervalID)
  }
}