import "./Header.css"
import { useState, useEffect } from "react"

function Header() {
  
  const [semester, setSemester] = useState("Sp")
  const [year, setYear] = useState(0)
  const [years, setYears] = useState([])

  useEffect( () => {
    /* Simulate Data */
    const tmpYears = []
    for (let i = 25; i >= 18; i--) {
      tmpYears.push(i)
    }
    setYears(tmpYears)
    /* ASSUME years always eixsts */
    setYear(years[0])
  }, [])

  return (
    <div className="header_main">
      <nav className="pages">
        <a>Courses</a>
        <a>My Courses</a>
      </nav>
      <div className="search_and_filter">
        <select value={semester} onChange={(event) => setSemester(event.target.value)}>
          <option value="F">Fall</option>
          <option value="Sp">Spring</option>
        </select>
        <select>
          {years.map((year, index) => (<option value={year}>{year}</option>))}
        </select>
        <input type="search" placeholder="Search Courses"/>
      </div>
    </div>
  );
}

export default Header;
