import React from "react"

// Icons
import { Search } from "lucide-react";

type SearchBarProps = {
    label?: string;
    value: string;
    onChange: (value: string) => void;   
}

const SearchBar = ({value, label = "Cari", onChange}: SearchBarProps)=>{
    return (
      <>
        <div className="relative">
          <input
            id="search-bar"
            name="search-bar"
            type="text"
            value={value}
            onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
              onChange(e.target.value)
            }
            className="
            rounded-sm w-full border-gray-200 
             focus:outline-none focus-visible:ring-2 focus-visible:ring-[#397e50]
            "
          />
          {value.length === 0 && (
            <label
              htmlFor="search-bar"
              className="absolute left-[15px] top-1/2 -translate-y-1/2
          text-slate-400 pointer-events-none flex items-center gap-2 "
            >
              <Search size={15} />
              {label}
            </label>
          )}
        </div>
      </>
    );
}

export default SearchBar;
