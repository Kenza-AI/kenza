import React from 'react'
import Filter from '../Filter.jsx'

export const header = (projectTitle, filterText, filterTextChangeHandler, filtersChangeHandler, filters) => (
    <header>
        <div>
            <h2 className='projects-header'><a href='/#/projects'>Projects</a> {projectTitle} </h2>
            <Filter text={filterText} placeholder='Search for training jobs by commit SHA' onChangeHandler={filterTextChangeHandler} />
            <ul className='jobs-type-filter'>
                <li>
                    <button onClick={() => filtersChangeHandler('')} className={`${filters.length > 0 ? null : 'active'}`}>All</button>
                </li>
                <li>
                    <button onClick={() => filtersChangeHandler('training')} className={`${filters.includes('training') ? 'active' : null}`}>Training</button>
                </li>
                <li>
                    <button onClick={() => filtersChangeHandler('tuning')} className={`${filters.includes('tuning') ? 'active' : null}`}>Hyperparameter Tuning</button>
                </li>
                {/* <li>
                    <button onClick={() => filtersChangeHandler('batchtransform')} className={`${filters.includes('batchtransform') ? 'active' : null}`}>Batch Transform</button>
                </li> */}
                <li>
                    <button onClick={() => filtersChangeHandler('endpoints')} className={`${filters.includes('endpoints') ? 'active' : null}`}>Jobs with Endpoints</button>
                </li>
            </ul>
        </div>
    </header>
)
