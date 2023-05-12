# prime service

## 时序图

```mermaid
sequenceDiagram
    autonumber
    participant client as client
    participant main as main thread
    participant work as work threads
    
    note over client: 整数N是否质数?
    client ->>+ main: 整数N是否质数
    note over main: 提取需要计算的整数，分批提交给工作线程计算
    loop 分批计算
        main -)+ work: 计算一批整数中的质数
        work --)- main: 计算完毕
        note over main: 将质数添加到质数数组（始终保持从小到大的顺序）
    end        
    main -->>- client: 计算完毕
    note over client: 输出整数N是否质数


```

## 参考资料

[mermaid 语法](https://cloud.tencent.com/developer/beta/article/1334691)

[Mermaid Live Editor](https://mermaid-js.github.io/mermaid-live-editor/edit)