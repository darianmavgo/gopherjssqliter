const App = () => {
    const gridRef = React.useRef(null);

    const columnDefs = [
        { field: 'id' },
        { field: 'name' },
        { field: 'email' },
        { field: 'age' },
    ];

    const onGridReady = (params) => {
        const datasource = {
            getRows: (params) => {
                console.log('Requesting rows:', params.startRow, 'to', params.endRow);
                
                // Call GopherJS exposed function
                // window.Backend.FetchRows(start, end, callback)
                if (window.Backend && window.Backend.FetchRows) {
                    window.Backend.FetchRows(params.startRow, params.endRow, (rows, lastRow) => {
                         params.successCallback(rows, lastRow);
                    });
                } else {
                    console.error("Backend not loaded");
                    params.failCallback();
                }
            }
        };

        params.api.setDatasource(datasource);
    };

    return (
        <div className="ag-theme-alpine">
            <agGridReact.AgGridReact
                ref={gridRef}
                columnDefs={columnDefs}
                rowModelType={'infinite'}
                onGridReady={onGridReady}
                cacheBlockSize={50}
                cacheOverflowSize={2}
                maxConcurrentDatasourceRequests={1}
                infiniteInitialRowCount={50}
                maxBlocksInCache={10}
            />
        </div>
    );
};

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(<App />);
