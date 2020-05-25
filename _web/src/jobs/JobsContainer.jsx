import { connect } from 'react-redux'

import Jobs from './Jobs.jsx'
import { fetchJobs, fetchLogs, compareJobs, deleteJobs, cancelJobs } from './api'
import { filterJobs, changeFilters, selectJobForComparison, clearJobsSelectedForComparison, dismissJobsComparison } from './actions'

const mapStateToProps = state => ({
   jobs: state.jobs,
   projectID: state.location.payload.projectID
})

const mapDispatchToProps = (dispatch, ownProps) => ({
   // Job fetching
   fetchJobsHandler: projectID => fetchJobs(dispatch, projectID, ownProps.config, ownProps.account),

   // Job deletion
   deleteJobsHandler: (projectID, jobIDs) => {
      if (confirm("Are you sure you want to delete this job?")) {
         deleteJobs(dispatch, projectID, jobIDs, ownProps.config, ownProps.account)
      }
   },

   // Job cancellation
   cancelJobsHandler: (projectID, jobIDs) => {
      if (confirm("Are you sure you want to cancel/stop this job?")) {
         cancelJobs(dispatch, projectID, jobIDs, ownProps.config, ownProps.account)
      }
   },

   // Job filtering
   filterTextChangeHandler: text => dispatch(filterJobs(text)),
   filtersChangeHandler: filter => dispatch(changeFilters(filter)),

   // Individual Job actions
   viewJobDetailsHandler: job => dispatch({ type: 'JOB', payload: { projectID: job.project.id, jobID: job.id } }),
   viewLogsHandler: (projectID, jobID) => fetchLogs(dispatch, projectID, jobID, ownProps.config, ownProps.account),

   // Job comparison
   clearJobsSelectedForComparisonHandler: () => dispatch(clearJobsSelectedForComparison()),
   compareJobsHandler: jobs => jobs.length > 1 ? compareJobs(dispatch, jobs[0].project.id, jobs, ownProps.config, ownProps.account) : window.alert('Please select at least two jobs.'),
   dismissJobsComparisonHandler: () => dispatch(dismissJobsComparison()),
   selectJobForComparisonHandler: (id, isSelectable, isSelectedForComparison) => isSelectable || isSelectedForComparison ? dispatch(selectJobForComparison(id)) :
      window.alert('Up to 3 FINISHED jobs can be selected for comparison.')
})

const JobsContainer = connect(mapStateToProps, mapDispatchToProps)(Jobs)
export default JobsContainer
