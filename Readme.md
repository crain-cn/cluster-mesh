# createapi
<pre>
mkdir temp && cd /temp
1. 初始化初始化项目；
2. 生产client。   -- 可以参照Hack中readme.md
3. devkubectl apply -f  deploy/crd/cluster_crd.yaml 
4. devkubectl apply -f  deploy/deploy.yaml 

测试
5. devkubectl apply -f  deploy/crd/cluster_cr.yaml 
</pre>
