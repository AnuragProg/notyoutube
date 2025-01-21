

## DAG Scheduling Algorithm

```
for {
    dag := <-dagChan // listen to dags
    go executeDAG(dag)
}

func executeDAG(dag){

    // create state
    createDAGState(dag, 'running')
    
    initialInput := dag.initialInput
    initialWorkers := getWorkersWithNoDependency(dag.workers)

    for _, worker := range initialWorkers {
        go executeWorker([input], worker) // executeWorker will update the worker state as well
    }
    
    for result := range workerCompletionChan {
        updateWorkerState(result.worker, 'completed') 
        
        deps := getDependenciesWhereWorkerIsSource(result.worker)
        
        for _, dep := range deps{
            allSourcesCompleted := All(Map(
                dep.srcs,
                func(src) bool {
                    return getWorkerState(src).status == 'completed'
                }
            ))
            
            inputs := getWorkerOutputs(dep.srcs)
            
            if allSourcesCompleted {
                go executeWorker(inputs, dep.tgt) // executeWorker will update the worker state as well
            }
        }
    }
}
```

